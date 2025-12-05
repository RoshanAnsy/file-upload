package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User model — one user can have many files
type User struct {
	ID        string      `gorm:"type:uuid;primaryKey"`
	Name      string      `gorm:"size:255;not null"`
	Email     string      `gorm:"uniqueIndex;size:255"`
	Password  string      `gorm:"size:255;not null"`
	Files     []FileStore `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ApiKey   []GenerateApiKey `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook — generate UUID automatically for each user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

// FileStore model — stores metadata about uploaded files
type FileStore struct {
	FileID      string    `gorm:"type:uuid;primaryKey"`
	UserID      string    `gorm:"type:uuid;not null;index"`
	FileName    string    `gorm:"size:255;not null"`
	FilePath    string    `gorm:"size:512;not null"`
	FileType    string    `gorm:"size:100"`  // e.g., "image/png", "application/pdf"
	FileSize    int64     // size in bytes
	Checksum    string    `gorm:"size:255"`  // e.g., SHA256 or MD5 hash
	IsPublic    bool      `gorm:"default:false"`
	Description string    `gorm:"size:512"`  // optional user-defined description
	User        User      `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BeforeCreate hook — generate UUID automatically for each file
func (f *FileStore) BeforeCreate(tx *gorm.DB) (err error) {
	f.FileID = uuid.New().String()
	return
}



type GenerateApiKey struct{
	Id string `gorm:"type:uuid;primaryKey"`
	UserID string `gorm:"type:uuid;not null;index"`
	KeyName string `gorm:"size:255;not null"`
	PublicKey string `gorm:"size:512;not null"`
	SecretKey string `gorm:"size:512;not null"`
	Status bool `gorm:"default:true"`
	User		User      `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (f *GenerateApiKey) BeforeCreate(tx *gorm.DB) (err error) {
	f.Id = uuid.New().String()
	return
}