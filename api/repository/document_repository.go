package repository

import (
	"context"
	"errors"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/models"
	"gorm.io/gorm"
)

type DocumentRepository interface {
    Create(ctx context.Context, doc *models.Document) error
    GetByID(ctx context.Context, id string) (*models.Document, error)
    GetByUserID(ctx context.Context, userID string) ([]*models.Document, error)
    UpdateStatus(ctx context.Context, id string, status models.DocStatus) error
    UpdateChunksCount(ctx context.Context, id string, count int) error
    Delete(ctx context.Context, id string) error
}

type documentRepository struct {
    db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
    return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, doc *models.Document) error {
    return r.db.WithContext(ctx).Create(doc).Error
}

func (r *documentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
    var doc models.Document
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &doc, err
}

func (r *documentRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Document, error) {
    var docs []*models.Document
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&docs).Error
    
    return docs, err
}

func (r *documentRepository) UpdateStatus(ctx context.Context, id string, status models.DocStatus) error {
    return r.db.WithContext(ctx).
        Model(&models.Document{}).
        Where("id = ?", id).
        Update("status", status).Error
}

func (r *documentRepository) UpdateChunksCount(ctx context.Context, id string, count int) error {
    return r.db.WithContext(ctx).
        Model(&models.Document{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "chunks_count": count,
            "status":       models.StatusReady,
        }).Error
}

func (r *documentRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&models.Document{}, "id = ?", id).Error
}