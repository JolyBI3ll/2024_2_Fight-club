package images

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioServiceInterface interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}

type MinioService struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioService(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (MinioServiceInterface, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket %s successfully created", bucketName)
	}

	return &MinioService{Client: client, BucketName: bucketName}, nil
}

func (m *MinioService) UploadFile(file *multipart.FileHeader, filePath string) (string, error) {
	fileObj, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileObj.Close()

	_, err = m.Client.PutObject(context.Background(), m.BucketName, filePath, fileObj, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}

	log.Printf("File %s successfully uploaded to %s", file.Filename, filePath)
	return filePath, nil
}

func (m *MinioService) DeleteFile(filePath string) error {
	err := m.Client.RemoveObject(context.Background(), m.BucketName, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	log.Printf("File %s successfully deleted", filePath)
	return nil
}
