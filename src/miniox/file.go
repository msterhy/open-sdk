package miniox

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/trancecho/open-sdk/libx"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

// 假设有一个全局的日志记录器
var logger = logrus.New()

func init() {
	// 配置日志轮转
	logFile := &lumberjack.Logger{
		Filename:   "./logs/download.log", // 日志文件路径
		MaxSize:    10,                    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 3,                     // 日志文件最多保存多少个备份
		MaxAge:     30,                    // 文件最多保存多少天
		Compress:   true,                  // 是否压缩
	}
	logger.SetOutput(logFile)

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// GetDownloadURL generates a presigned download URL for the specified object in MinIO.
func GetDownloadURL(c *gin.Context, bucketName string, objectName string, userfilename string) (string, error) {
	logger := logger.WithFields(logrus.Fields{
		"bucketName":   bucketName,
		"objectName":   objectName,
		"userfilename": userfilename,
	})

	logger.Info("GetDownloadURL function started")

	if bucketName == "" || objectName == "" {
		err := errors.New("Bucket and object parameters are required")
		logger.Error(err)
		return "", err
	}

	// 使用 PresignedDownload 函数生成预签名下载链接
	presignedURL, err := PresignedDownload(c, bucketName, objectName)
	if err != nil {
		logger.WithError(err).Error("Failed to generate presigned URL")
		return "", err
	}

	logger.Info("Presigned URL generated successfully")
	return presignedURL, nil
}

func Upload(c *gin.Context, bucketName string, objectName string) {
	if bucketName == "" || objectName == "" {
		libx.Err(c, http.StatusInternalServerError, "Bucket and object parameters are required", nil)
		return
	}

	file, err := c.FormFile("uploadFile")
	if err != nil {
		libx.Err(c, http.StatusBadRequest, "Error retrieving the file", err)
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		libx.Err(c, http.StatusInternalServerError, "Error opening the file", err)
		return
	}
	defer func(srcFile multipart.File) {
		err := srcFile.Close()
		if err != nil {
			libx.Err(c, http.StatusInternalServerError, "Error closing the file", err)
		}
	}(srcFile)

	_, err = MinioClient.PutObject(context.Background(), bucketName, objectName, srcFile, file.Size, minio.PutObjectOptions{})
	if err != nil {
		libx.Err(c, http.StatusInternalServerError, "Error uploading the file", err)
		return
	}
	libx.Ok(c, "File uploaded successfully", nil)
}
func DownloadToLocal(c *gin.Context, client *minio.Client, bucketName string, objectName string, filePath string) error {
	// 构建保存路径，将文件存储到当前目录的 tmp 文件夹下
	tmpFilePath := fmt.Sprintf("tmp/%s", filePath)

	// 从MinIO存储桶中下载文件
	err := client.FGetObject(c, bucketName, objectName, tmpFilePath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error downloading file from MinIO: %v", err)
		libx.Err(c, http.StatusInternalServerError, "从MinIO下载文件失败", err)
		return err
	}
	log.Println("File downloaded from MinIO successfully")
	return nil
}

func AddFolder(bucketName string, folderName string) error {
	if bucketName == "" || folderName == "" {
		return errors.New("bucket and folder parameters are required")
	}
	_, err := MinioClient.PutObject(context.Background(), bucketName, folderName, nil, 0, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func RenameFolder(bucketName, oldFolderPath, newFolderPath string) error {
	if bucketName == "" || oldFolderPath == "" || newFolderPath == "" {
		return errors.New("bucket, old folder and new folder parameters are required")
	}
	log.Println("Renaming folder from", oldFolderPath, "to", newFolderPath)

	//// 确保路径以斜杠结尾
	//if !strings.HasSuffix(oldFolderPath, "/") {
	//	oldFolderPath += "/"
	//}
	//if !strings.HasSuffix(newFolderPath, "/") {
	//	newFolderPath += "/"
	//}

	ctx := context.Background()

	// 列出旧文件夹下的所有对象
	log.Println("Listing objects in bucket:", bucketName, "with prefix:", oldFolderPath)
	objectCh := MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    oldFolderPath,
		Recursive: true,
	})

	// 如果对象列表为空，则返回错误,说明没有发生更新。
	if objectCh == nil {
		log.Println("No objects found in the folder:", oldFolderPath)
		return errors.New("no objects found in the folder: " + oldFolderPath)
	}

	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing object:", object.Err)
			return object.Err
		}

		oldObjectKey := object.Key
		//log.Println("Found object:", oldObjectKey)
		newObjectKey := newFolderPath + strings.TrimPrefix(oldObjectKey, oldFolderPath)
		//log.Println("New object key:", newObjectKey)

		// 检查源对象是否存在
		log.Println("Checking if source object exists:", oldObjectKey)
		_, err := MinioClient.StatObject(ctx, bucketName, oldObjectKey, minio.StatObjectOptions{})
		if err != nil {
			log.Println("Source object does not exist:", oldObjectKey)
			return errors.New("source object does not exist: " + oldObjectKey)
		}

		// 复制对象到新位置
		log.Printf("Copying %s to %s", oldObjectKey, newObjectKey)
		_, err = MinioClient.CopyObject(ctx, minio.CopyDestOptions{
			Bucket: bucketName,
			Object: newObjectKey,
		}, minio.CopySrcOptions{
			Bucket: bucketName,
			Object: oldObjectKey,
		})
		if err != nil {
			log.Println("Error copying object:", err)
			return err
		}

		// 删除旧对象
		log.Printf("Removing old object: %s", oldObjectKey)
		err = MinioClient.RemoveObject(ctx, bucketName, oldObjectKey, minio.RemoveObjectOptions{})
		if err != nil {
			log.Println("Error removing old object:", err)
			return err
		}
	}
	log.Println("Folder rename completed successfully")

	return nil
}

func RemoveFolder(bucketName, folderPath string) error {
	if bucketName == "" || folderPath == "" {
		return errors.New("bucket and folder parameters are required")
	}

	//// 确保路径以斜杠结尾
	//if !strings.HasSuffix(folderPath, "/") {
	//	folderPath += "/"
	//}

	ctx := context.Background()

	// 列出文件夹下的所有对象
	log.Println("Listing objects in bucket:", bucketName, "with prefix:", folderPath)
	objectCh := MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    folderPath,
		Recursive: true,
	})

	// 如果对象列表为空，则返回错误
	if objectCh == nil {
		log.Println("No objects found in the folder:", folderPath)
		return errors.New("no objects found in the folder: " + folderPath)
	}

	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing object:", object.Err)
			return object.Err
		}

		objectKey := object.Key
		//log.Println("Found object:", objectKey)

		// 删除对象
		log.Printf("Removing object: %s", objectKey)
		err := MinioClient.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
		if err != nil {
			log.Println("Error removing object:", err)
			return err
		}
	}
	log.Println("Folder removed successfully")

	return nil
}

func AddFile(bucketName, objectName string, file multipart.File, size int64, contentType string) error {
	if bucketName == "" || objectName == "" {
		return errors.New("bucket and object parameters are required")
	}

	_, err := MinioClient.PutObject(context.Background(), bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}
	return nil
}
