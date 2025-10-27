package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/repository"
)

var s3Client *s3.Client

func Init() {
    awsEndpoint := os.Getenv("AWS_ENDPOINT")     // http://localstack:4566
    awsRegion := os.Getenv("AWS_REGION")         // eu-central-1
    
    slog.Info("initializing AWS S3 client",
        "endpoint", awsEndpoint,
        "region", awsRegion,
        "is_localstack", awsEndpoint != "")

    ctx := context.TODO()
    var cfg aws.Config
    var err error

    if awsEndpoint != "" {
        slog.Info("using LocalStack endpoint")
        
        cfg, err = config.LoadDefaultConfig(ctx,
            config.WithRegion(awsRegion),
            config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
                "test",  // Access Key
                "test",  // Secret Key
                "",      // Session Token
            )),
        )
        
        if err != nil {
            slog.Error("failed to load AWS config", "error", err)
            panic(err)
        }

        slog.Info("creating s3 client")
        s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
            o.BaseEndpoint = aws.String(awsEndpoint)
            o.UsePathStyle = true  // LocalStack
        })
        
    } else {
        slog.Info("using real AWS endpoint")
    }

    slog.Info("S3 client initialized successfully")
}

func GenerateS3Key(userID, fileName string) string {
    key := fmt.Sprintf("documents/%s/%s-%s", userID, uuid.New().String(), fileName)
    slog.Debug("generated S3 key", "key", key, "user_id", userID)
    return key
}

func GeneratePresignedUploadURL(key, contentType string, expiration time.Duration) (string, error) {
    bucketName := os.Getenv("AWS_S3_BUCKET")
    
    slog.Info("generating presigned upload URL",
        "bucket", bucketName,
        "key", key,
        "content_type", contentType,
        "expiration", expiration)

    presignClient := s3.NewPresignClient(s3Client)
    
    request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
        Bucket:      &bucketName,
        Key:         &key,
        ContentType: &contentType,
    }, func(opts *s3.PresignOptions) {
        opts.Expires = expiration
    })

    if err != nil {
        slog.Error("failed to generate presigned URL",
            "error", err,
            "bucket", bucketName,
            "key", key)
        return "", err
    }

    slog.Debug("presigned URL generated", "url", request.URL)


    // Client isn't aware of localstack - it knows only for localhost
    presignedUrl := strings.Replace(request.URL, "http://localstack:4566", "http://localhost:4566", 1)
    return presignedUrl, nil
}

func CreateDocumentRecord(userID, fileName, s3Key string, fileSize int64) (string, error) {
    documentID := uuid.New().String()
    
    slog.Info("creating document record",
        "document_id", documentID,
        "user_id", userID,
        "file_name", fileName,
        "file_size", fileSize)

    query := `
        INSERT INTO documents (id, user_id, file_name, s3_key, file_size, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
    `
    
    _, err := repository.DB.Exec(query, documentID, userID, fileName, s3Key, fileSize, "uploading")
    if err != nil {
        slog.Error("failed to create document record",
            "error", err,
            "document_id", documentID,
            "user_id", userID)
        return "", err
    }

    slog.Info("document record created successfully", "document_id", documentID)
    return documentID, nil
}

func SaveDocumentRecord(documentID string) error {
    query := `
        UPDATE documents
        SET status = $1,
            uploaded_at = NOW()
        WHERE id = $2
    `

    _, err := repository.DB.Exec(query, "uploaded", documentID)
    if err != nil {
        slog.Error("failed to update document record",
            "error", err,
            "document_id", documentID)
        return err
    }

    slog.Info("document record updated successfully", "document_id", documentID)
    return nil
}

