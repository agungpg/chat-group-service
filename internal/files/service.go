package files

import (
	"context"
	"fmt"
	"time"

	"github.com/agungpg/group-chat-service/config"
	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type Service struct {
	s3      *s3.Client
	repo    *Repository
	strConf *config.StorageConfig
}

func NewService(repo *Repository, s3 *s3.Client) *Service {
	return &Service{
		repo: repo,
		s3:   s3,
	}
}

func PresignUploadUrl(ctx context.Context, presigner *s3.PresignClient, bucket, key, contentType string, expires time.Duration) (string, error) {
	in := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType), // client MUST send the same Content-Type header
	}

	out, err := presigner.PresignPutObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func PresignViewURL(
	ctx context.Context,
	presigner *s3.PresignClient,
	bucket, key string,
	expires time.Duration,
) (string, error) {

	in := &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("inline"),
	}

	out, err := presigner.PresignGetObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func PresignDownloadURL(
	ctx context.Context,
	presigner *s3.PresignClient,
	bucket, key, filename string,
	expires time.Duration,
) (string, error) {

	disposition := fmt.Sprintf(`attachment; filename="%s"`, filename)

	in := &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(disposition),
	}

	out, err := presigner.PresignGetObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
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
		bucket = s.strConf.PrivateBucket
		scope = FileScopeAvatar
		key = string(scope) + "/" + fileName
	} else {
		bucket = s.strConf.PublicBucket
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
		CreatedBy:        userId,
	}

	err := s.repo.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	presigner := s3.NewPresignClient(s.s3)
	expiresIn := 15 * time.Minute
	url, err := PresignUploadUrl(ctx, presigner, bucket, key, f.ContentType, expiresIn)
	if err != nil {
		return nil, err
	}

	return &UploadPresignResponseDTO{
		PresignedUrl: url,
		ExpiresIn:    expiresIn,
		FileId:       fileId,
	}, nil
}
