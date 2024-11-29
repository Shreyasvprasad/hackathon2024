package fileupload

import "github.com/minio/minio-go/v7"

func UploadFile(bucketName, filePath, filename string) (string, error) {
	client, err := minio.New("minio.server.com", "accessKey", "secretKey", false)
	if err != nil {
		return "", err
	}

	// Upload the file
	_, err = client.FPutObject(bucketName, filename, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// Return the file URL
	return "http://minio.server.com/" + bucketName + "/" + filename, nil
}
