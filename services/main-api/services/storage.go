package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/config"
)

// UploadFile uploads a file to the specified MinIO bucket
func UploadFile(bucketName string, objectName string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ctx := context.Background()

	// Upload the file to specific bucket
	info, err := config.MinioClient.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// Format URL assuming public read bucket
	// Format: http://localhost:9000/{bucketName}/{objectName}
	fileURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, info.Key)
	return fileURL, nil
}

// GetPresignedURL returns a temporary download URL for private documents (e.g. invoices)
func GetPresignedURL(bucketName string, objectName string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	
	reqParams := make(url.Values)
	// Example: reqParams.Set("response-content-disposition", "attachment; filename=\"invoice.pdf\"")

	presignedURL, err := config.MinioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// DeleteFile removes a file from MinIO (e.g. when a project is deleted)
func DeleteFile(bucketName string, objectName string) error {
	ctx := context.Background()
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err := config.MinioClient.RemoveObject(ctx, bucketName, objectName, opts)
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

// EnsureBucketsExist makes sure the required buckets exist (if docker minio-init failed for some reason)
func EnsureBucketsExist() {
	ctx := context.Background()
	buckets := []string{"cad-parts", "renders", "documents"}

	for _, bucketName := range buckets {
		exists, err := config.MinioClient.BucketExists(ctx, bucketName)
		if err != nil {
			log.Printf("Error checking bucket %s: %v\n", bucketName, err)
			continue
		}
		
		if !exists {
			err = config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Printf("Error creating bucket %s: %v\n", bucketName, err)
			} else {
				log.Printf("Successfully created bucket %s\n", bucketName)
			}
		}
	}
}
