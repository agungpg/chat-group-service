package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Signer struct {
	presigner *s3.PresignClient
}

func NewSigner(presigner *s3.PresignClient) *Signer {
	return &Signer{
		presigner: presigner,
	}
}

func (s *Signer) GetPresignUploadUrl(ctx context.Context, bucket, key, contentType string, expires time.Duration) (string, error) {
	in := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType), // client MUST send the same Content-Type header
	}

	out, err := s.presigner.PresignPutObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func (s *Signer) GetPresignViewURL(
	ctx context.Context,
	bucket, key string,
	expires time.Duration,
) (string, error) {

	in := &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("inline"),
	}

	out, err := s.presigner.PresignGetObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}

func (s *Signer) GetPresignDownloadURL(
	ctx context.Context,
	bucket, key, filename string,
	expires time.Duration,
) (string, error) {

	disposition := fmt.Sprintf(`attachment; filename="%s"`, filename)

	in := &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(disposition),
	}

	out, err := s.presigner.PresignGetObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = expires
	})
	if err != nil {
		return "", err
	}

	return out.URL, nil
}
