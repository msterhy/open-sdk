package miniox

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/trancecho/open-sdk/config"
	"log"
	"time"
)

var MinioClient *minio.Client

func MinioInit() (ok bool) {
	endpoint, accessKeyID, secretAccessKey, useSSL := MinioProfile()
	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Println("MinIO client configuration not set")
		return false
	}

	//log.Println("MinIO client configuration set")
	var err error
	// 初始化MinIO客户端
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Error initializing MinIO client: %v", err)
		return false
	}
	log.Println("MinIO client initialized successfully")

	return true
}

func MinioProfile() (string, string, string, bool) {
	minioConfig := config.GetConfig().Minio
	// MinIO配置
	endpoint := minioConfig.Endpoint
	accessKeyID := minioConfig.AccessKeyID
	secretAccessKey := minioConfig.SecretAccessKey
	useSSL := minioConfig.UseSSL
	return endpoint, accessKeyID, secretAccessKey, useSSL
}

// 生成预签名下载链接 - 有效期24小时 超时时间5秒 支持取消
func PresignedDownload(c *gin.Context, bucketName string, objectName string) (string, error) {
	if bucketName == "" || objectName == "" {
		return "", errors.New("bucketName 和 objectName 是必需的")
	}

	// 使用 Gin 的上下文，支持取消
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 调用 MinIO SDK 生成预签名链接
	expiry := 24 * time.Hour // 24 hours in seconds
	presignedURL, err := MinioClient.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		// 检查是否是取消操作导致的错误
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("请求超时: %v", err)
			return "", errors.New("生成预签名链接超时，请稍后重试")
		}
		if ctx.Err() == context.Canceled {
			log.Printf("请求被取消: %v", err)
			return "", errors.New("操作已取消")
		}
		log.Printf("生成预签名链接失败: %v", err)
		return "", errors.New("生成预签名链接失败，请检查存储桶和对象名称")
	}

	return presignedURL.String(), nil

}
