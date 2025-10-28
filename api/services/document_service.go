package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/models"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/repository"
)

type DocumentService struct {
    repo      repository.DocumentRepository
    s3Service *S3Service
}

func NewDocumentService(repo repository.DocumentRepository, s3 *S3Service) *DocumentService {
    return &DocumentService{
        repo:      repo,
        s3Service: s3,
    }
}

func (s *DocumentService) CreateDocument(ctx context.Context, userID, fileName, fileType string, fileSize int64) (*models.Document, string, error) {
    documentID := uuid.New()
    s3Key := s.s3Service.GenerateS3Key(userID, fileName)
    
    slog.Info("creating document",
        "document_id", documentID,
        "user_id", userID,
        "file_name", fileName,
        "file_size", fileSize)
    
    doc := &models.Document{
        ID:        documentID,
        UserID:    userID,
        FileName:  fileName,
        S3Key:     s3Key,
        FileSize:  fileSize,
        Status:    models.StatusUploading,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, doc); err != nil {
        slog.Error("failed to create document", "error", err)
        return nil, "", err
    }
    
    uploadURL, err := s.s3Service.GeneratePresignedUploadURL(s3Key, fileType, 15*time.Minute)
    if err != nil {
        slog.Error("failed to generate presigned URL", "error", err)
        return nil, "", err
    }
    
    slog.Info("document created successfully", "document_id", documentID)
    return doc, uploadURL, nil
}

func (s *DocumentService) MarkAsUploaded(ctx context.Context, documentID string) error {
    slog.Info("marking document as uploaded", "document_id", documentID)
    return s.repo.UpdateStatus(ctx, documentID, models.StatusReady)
}

func (s *DocumentService) GetUserDocuments(ctx context.Context, userID string) ([]*models.Document, error) {
    return s.repo.GetByUserID(ctx, userID)
}

func (s *DocumentService) GetDocumentByID(ctx context.Context, documentID string) (*models.Document, error) {
    return s.repo.GetByID(ctx, documentID)
}