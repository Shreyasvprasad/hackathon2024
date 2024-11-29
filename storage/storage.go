package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinIO() {
	client, err := minio.New("play.min.io", &minio.Options{
		Creds:  credentials.NewStaticV4("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY", ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	minioClient = client
}

func UploadFile(bucketName, objectName, filePath string) (string, error) {
	_, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://play.min.io/%s/%s", bucketName, objectName), nil
}
