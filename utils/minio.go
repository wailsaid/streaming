package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinioClient() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %v", err)
	}

	minioClient = client

	// Create buckets if they don't exist
	buckets := []string{"videos", "thumbnails"}
	ctx := context.Background()

	for _, bucket := range buckets {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("failed to check bucket existence: %v", err)
		}

		if !exists {
			err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("failed to create bucket %s: %v", bucket, err)
			}
		}
	}

	return nil
}

func GetMinioClient() *minio.Client {
	return minioClient
} 