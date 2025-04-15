package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client = *minio.Client

type Config struct {
	Address  string `yaml:"address" env:"MINIO_ADDRESS"`
	User     string `yaml:"user" env:"MINIO_USER"`
	Password string `env:"MINIO_PASSWORD"`
	Bucket   string `yaml:"bucket" env:"MINIO_BUCKET"`
	UseSSL   bool   `yaml:"use_ssl" env:"MINIO_USE_SSL"`
}

func New(ctx context.Context, cfg *Config) (Client, error) {
	client, err := minio.New(
		cfg.Address,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.User, cfg.Password, "",
			),
			Secure: cfg.UseSSL,
		},
	)

	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
			 {
			  "Effect": "Allow",
			  "Principal": {"AWS": ["*"]},
			  "Action": ["s3:GetObject"],
			  "Resource": ["arn:aws:s3:::%s/*"]
			 }
			]
			}`, cfg.Bucket)

		client.SetBucketPolicy(ctx, cfg.Bucket, policy)
	}

	return client, nil
}
