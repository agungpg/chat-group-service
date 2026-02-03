package files

import (
	"context"
	"time"

	"github.com/agungpg/group-chat-service/internal/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageAdapter interface {
	PresignUploadURL(ctx context.Context, bucket, key, contentType string, expires time.Duration) (string, error)
	PresignViewURL(ctx context.Context, bucket, key string, expires time.Duration) (string, error)
	PresignDownloadURL(ctx context.Context, bucket, key, filename string, expires time.Duration) (string, error)
	Exists(ctx context.Context, bucket, key string) error
}

type S3Adapter struct {
	signer *storage.Signer
	object *storage.Object
}

func NewS3Adapter(s3Client *s3.Client) *S3Adapter {
	presigner := s3.NewPresignClient(s3Client)

	return &S3Adapter{
		signer: storage.NewSigner(presigner),
		object: storage.NewObject(s3Client),
	}
}

func (a *S3Adapter) PresignUploadURL(ctx context.Context, bucket, key, contentType string, expires time.Duration) (string, error) {
	return a.signer.GetPresignUploadUrl(ctx, bucket, key, contentType, expires)
}

func (a *S3Adapter) PresignViewURL(ctx context.Context, bucket, key string, expires time.Duration) (string, error) {
	return a.signer.GetPresignViewURL(ctx, bucket, key, expires)
}

func (a *S3Adapter) PresignDownloadURL(ctx context.Context, bucket, key, filename string, expires time.Duration) (string, error) {
	return a.signer.GetPresignDownloadURL(ctx, bucket, key, filename, expires)
}

func (a *S3Adapter) Exists(ctx context.Context, bucket, key string) error {
	_, err := a.object.Exist(ctx, bucket, key)
	return err
}
