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
)

type S3Service struct {
	client     *s3.Client
	bucketName string
}

func InitS3() (*S3Service, error) {
	awsEndpoint := os.Getenv("AWS_ENDPOINT") 
	awsRegion := os.Getenv("AWS_REGION")     
	bucketName := os.Getenv("AWS_S3_BUCKET")

	slog.Info("initializing AWS S3 client",
		"endpoint", awsEndpoint,
		"region", awsRegion,
		"bucket", bucketName,
		"is_localstack", awsEndpoint != "")

	ctx := context.TODO()
	var cfg aws.Config
	var err error

	if awsEndpoint != "" {
		slog.Info("using LocalStack endpoint")

		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				"test", // Access Key
				"test", // Secret Key
				"",     // Session Token
			)),
		)

		if err != nil {
			slog.Error("failed to load AWS config", "error", err)
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}

		client := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(awsEndpoint)
			o.UsePathStyle = true // LocalStack zahteva path-style
		})

		slog.Info("S3 client initialized successfully (LocalStack)")
		return &S3Service{
			client:     client,
			bucketName: bucketName,
		}, nil

	} else {
		slog.Info("using real AWS endpoint")

		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(awsRegion),
		)

		if err != nil {
			slog.Error("failed to load AWS config", "error", err)
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}

		client := s3.NewFromConfig(cfg)

		slog.Info("S3 client initialized successfully (AWS)")
		return &S3Service{
			client:     client,
			bucketName: bucketName,
		}, nil
	}
}


func (s *S3Service) GenerateS3Key(userID, fileName string) string {
	key := fmt.Sprintf("documents/%s/%s-%s", userID, uuid.New().String(), fileName)
	slog.Debug("generated S3 key", "key", key, "user_id", userID)
	return key
}


func (s *S3Service) GeneratePresignedUploadURL(key, contentType string, expiration time.Duration) (string, error) {
	slog.Info("generating presigned upload URL",
		"bucket", s.bucketName,
		"key", key,
		"content_type", contentType,
		"expiration", expiration)

	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		slog.Error("failed to generate presigned URL",
			"error", err,
			"bucket", s.bucketName,
			"key", key)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	slog.Debug("presigned URL generated", "url", request.URL)

	// Client isn't aware of localstack - it knows only for localhost
	presignedURL := strings.Replace(request.URL, "http://localstack:4566", "http://localhost:4566", 1)
	return presignedURL, nil
}