package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Object struct {
	client *s3.Client
}

func NewObject(client *s3.Client) *Object {
	return &Object{
		client: client,
	}
}

func (o *Object) Exist(ctx context.Context, bucket, key string) (bool, error) {
	headObjectInput := &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	_, err := o.client.HeadObject(ctx, headObjectInput)
	if err != nil {
		return false, err
	}
	return true, nil
}
