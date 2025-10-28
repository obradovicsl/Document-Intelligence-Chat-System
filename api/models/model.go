package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
    ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
    UserID      string         `gorm:"type:varchar(255);not null;index:idx_user_id" json:"userId"`
    Question    string         `gorm:"type:text;not null" json:"question"`
    Answer       string         `gorm:"type:text;not null" json:"answer"`
    AskedAt   time.Time      `gorm:"autoCreateTime" json:"askedAt"`
}

type Document struct {
    ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
    UserID      string         `gorm:"type:varchar(255);not null;index:idx_user_id" json:"userId"`
    FileName    string         `gorm:"type:varchar(500);not null" json:"fileName"`
    S3Key       string         `gorm:"type:varchar(500);not null" json:"s3Key"`
    FileSize    int64          `gorm:"not null" json:"fileSize"`
    Status      DocStatus      `gorm:"type:varchar(50);not null;index:idx_user_status,priority:2" json:"status"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type DocStatus string

const (
    StatusUploading  DocStatus = "uploading"
    StatusReady      DocStatus = "ready"
)