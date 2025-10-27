package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/repository"
)

var s3Client *s3.Client

func init() {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(os.Getenv("AWS_REGION")),
    )
    if err != nil {
        panic(err)
    }
    s3Client = s3.NewFromConfig(cfg)
}

func GenerateS3Key(userID, fileName string) string {
    return fmt.Sprintf("documents/%s/%s-%s", userID, uuid.New().String(), fileName)
}

func GeneratePresignedUploadURL(key, contentType string, expiration time.Duration) (string, error) {
    bucketName := os.Getenv("AWS_S3_BUCKET")

    presignClient := s3.NewPresignClient(s3Client)

    request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
        Bucket:      &bucketName,
        Key:         &key,
        ContentType: &contentType,
    }, func(opts *s3.PresignOptions) {
        opts.Expires = expiration
    })

    if err != nil {
        return "", err
    }

    return request.URL, nil
}

func CreateDocumentRecord(userID, fileName, s3Key string, fileSize int64) (string, error) {
    documentID := uuid.New().String()

    query := `
        INSERT INTO documents (id, user_id, file_name, s3_key, file_size, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
    `

    _, err := repository.DB.Exec(query, documentID, userID, fileName, s3Key, fileSize, "uploading")
    if err != nil {
        return "", err
    }

    return documentID, nil
}