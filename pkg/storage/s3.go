package storage

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/agungpg/group-chat-service/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/config"
)

var S3Client *s3.Client

type Config struct {
	AccountID string
	AccessKey string
	SecretKey string
}

func NewConfig() *Config {
	return &Config{
		AccountID: utils.GetEnvOrDefault("STORAGE_ACCOUNT_ID", ""),
		AccessKey: utils.GetEnvOrDefault("STORAGE_ACCESS_KEY", ""),
		SecretKey: utils.GetEnvOrDefault("STORAGE_SECRET_KEY", ""),
	}
}

func (c *Config) GetEndPoint() string {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID)

	return endpoint
}

func (c *Config) GetDefaultConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, "")),
	)
}

func (c *Config) NewConnection(endpoint string, awsConf aws.Config, httpClient *http.Client) {

	S3Client = s3.NewFromConfig(awsConf, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.HTTPClient = httpClient
	})

}

func Connect(ctx context.Context) error {
	c := NewConfig()
	endpoint := c.GetEndPoint()

	awsConf, err := c.GetDefaultConfig(ctx)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}

	c.NewConnection(endpoint, awsConf, httpClient)

	return nil
}
