package files

import (
	"context"
	"fmt"
	"time"

	"github.com/agungpg/group-chat-service/config"
	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	storage  StorageAdapter
	strgConf *config.StorageConfig
}

func NewService(repo *Repository, storage StorageAdapter, strgConf *config.StorageConfig) *Service {
	return &Service{
		repo:     repo,
		storage:  storage,
		strgConf: strgConf,
	}
}

func (s *Service) Exist(ctx context.Context, bucket, key string) error {
	return s.storage.Exists(ctx, bucket, key)
}

func (s *Service) CreatePresignUpload(ctx context.Context, f UploadPresignRequestDTO, userId string) (*UploadPresignResponseDTO, error) {

	var (
		bucket string
		key    string
		scope  FileScope
	)
	fileId := uuid.New().String()
	fileName := fileId + utils.ExtFromContentType(f.ContentType)

	if f.FileType == "avatar" {
		bucket = s.strgConf.PublicBucket
		scope = FileScopeAvatar
		key = string(scope) + "/" + fileName
	} else {
		bucket = s.strgConf.PrivateBucket
		scope = FileScopeChatAttachment
		key = string(scope) + "/" + fileName
	}
	file := &Files{
		ID:               fileId,
		Bucket:           bucket,
		ObjectKey:        key,
		OriginalFilename: f.FileName,
		SizeBytes:        f.SizeBytes,
		Scope:            scope,
		Status:           FileScope(FileStatusPending),
		CreatedBy:        userId,
		CreatedAt:        time.Now(),
	}

	err := s.repo.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	expiresIn := 15 * time.Minute
	url, err := s.storage.PresignUploadURL(ctx, bucket, key, f.ContentType, expiresIn)
	if err != nil {
		return nil, err
	}

	return &UploadPresignResponseDTO{
		PresignedUrl: url,
		ExpiresIn:    expiresIn,
		FileId:       fileId,
	}, nil
}

func (s *Service) GetPresignView(ctx context.Context, fileId string) (*UploadPresignResponseDTO, error) {
	file, err := s.repo.GetFileById(ctx, fileId)
	if err != nil {
		return nil, err
	}

	expiresIn := 6 * time.Hour

	url, err := s.storage.PresignViewURL(ctx, file.Bucket, file.ObjectKey, expiresIn)

	return &UploadPresignResponseDTO{
		PresignedUrl: url,
		ExpiresIn:    expiresIn,
		FileId:       fileId,
	}, nil
}

func (s *Service) UpdateFileStatus(ctx context.Context, fileId string, status FileStatus) error {
	file, err := s.repo.GetFileById(ctx, fileId)
	if err != nil {
		return err
	}

	if file == nil {
		return fmt.Errorf("File not found!")
	}

	if status == FileStatusReady {
		err = s.Exist(ctx, file.Bucket, file.ObjectKey)

		if err != nil {
			return err
		}
	}

	err = s.repo.UpdateFileStatus(ctx, fileId, status)
	return err
}
