package s3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BucketName = "your-bucket-name"
	Region     = "us-east-1"
)

var s3Client *s3.Client

func main() {
	// Загружаем конфиг AWS
	//cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(Region))
	//if err != nil {
	//	log.Fatalf("Ошибка загрузки AWS конфигурации: %v", err)
	//}
	//
	//// Создаем клиент S3
	//s3Client = s3.NewFromConfig(cfg)
	//
	//r := gin.Default()
	//
	//// Генерация presigned URL
	//r.=("/presigned-url", func(c *gin.Context) {
	//	key := c.Query("key")       // Имя файла в S3
	//	action := c.Query("action") // "upload" или "download"
	//
	//	if key == "" || (action != "upload" && action != "download") {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
	//		return
	//	}
	//
	//	url, err := generatePresignedURL(key, action)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
	//		return
	//	}
	//
	//	c.JSON(http.StatusOK, gin.H{"url": url})
	//})
	//
	//r.Run(":8080")
}

// Генерируем presigned URL
func generatePresignedURL(key string, action string) (string, error) {
	ctx := context.TODO()

	var presignClient *s3.PresignClient
	presignClient = s3.NewPresignClient(s3Client)

	expiration := 15 * time.Minute

	var presignedReq *s3.PresignedPostRequest
	var err error

	switch action {
	case "upload":
		presignedReq, err = presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(BucketName),
			Key:    aws.String(key),
		}, s3.WithPresignExpires(expiration))
	case "download":
		presignedReq, err = presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(BucketName),
			Key:    aws.String(key),
		}, s3.WithPresignExpires(expiration))
	}

	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return presignedReq.URL, nil
}
