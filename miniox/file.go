package miniox

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/trancecho/open-sdk/libx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func Download(c *gin.Context, bucketName string, objectName string) {
	if bucketName == "" || objectName == "" {
		libx.Err(c, http.StatusInternalServerError, "Bucket and object parameters are required", libx.ErrOptions{})
		return
	}

	object, err := MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		libx.Err(c, http.StatusInternalServerError, "Could not retrieve the file", libx.ErrOptions{})
		return
	}
	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			libx.Err(c, http.StatusInternalServerError, "Error closing the file", libx.ErrOptions{})
		}
	}(object)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectName))
	c.Header("Content-Type", "application/octet-stream")
	if _, err = io.Copy(c.Writer, object); err != nil {
		libx.Err(c, http.StatusInternalServerError, "Error sending the file", libx.ErrOptions{})
		return
	}
	libx.Ok(c, "File sent successfully", nil)
}

func Upload(c *gin.Context, bucketName string, objectName string) {
	if bucketName == "" || objectName == "" {
		libx.Err(c, http.StatusInternalServerError, "Bucket and object parameters are required", libx.ErrOptions{})
		return
	}

	file, err := c.FormFile("uploadFile")
	if err != nil {
		libx.Err(c, http.StatusBadRequest, "Error retrieving the file", libx.ErrOptions{})
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		libx.Err(c, http.StatusInternalServerError, "Error opening the file", libx.ErrOptions{})
		return
	}
	defer func(srcFile multipart.File) {
		err := srcFile.Close()
		if err != nil {
			libx.Err(c, http.StatusInternalServerError, "Error closing the file", libx.ErrOptions{})
		}
	}(srcFile)

	_, err = MinioClient.PutObject(context.Background(), bucketName, objectName, srcFile, file.Size, minio.PutObjectOptions{})
	if err != nil {
		libx.Err(c, http.StatusInternalServerError, "Error uploading the file", libx.ErrOptions{})
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
		libx.Err(c, http.StatusInternalServerError, "从MinIO下载文件失败", libx.ErrOptions{})
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
