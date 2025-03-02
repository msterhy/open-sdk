package miniox

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type MinioClientInterface interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error)
	PutObject(ctx context.Context, bucketName, objectName string, data io.Reader, length int64, opts minio.PutObjectOptions) (int64, error)
	FGetObject(ctx context.Context, bucketName, objectName string, filePath string, opts minio.GetObjectOptions) error
}

// 为测试提供的 Mock Client
type MockMinioClient struct {
	GetObjectFunc  func(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error)
	PutObjectFunc  func(ctx context.Context, bucketName, objectName string, data io.Reader, length int64, opts minio.PutObjectOptions) (int64, error)
	FGetObjectFunc func(ctx context.Context, bucketName, objectName string, filePath string, opts minio.GetObjectOptions) error
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (io.ReadCloser, error) {
	return m.GetObjectFunc(ctx, bucketName, objectName, opts)
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, data io.Reader, length int64, opts minio.PutObjectOptions) (int64, error) {
	return m.PutObjectFunc(ctx, bucketName, objectName, data, length, opts)
}

func (m *MockMinioClient) FGetObject(ctx context.Context, bucketName, objectName string, filePath string, opts minio.GetObjectOptions) error {
	return m.FGetObjectFunc(ctx, bucketName, objectName, filePath, opts)
}
