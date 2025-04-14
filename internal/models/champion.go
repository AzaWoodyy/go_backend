package models

import (
	"time"

	"gorm.io/gorm"
)

type Champion struct {
	gorm.Model
	ID        string    `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Blurb     string    `gorm:"type:text" json:"blurb"`
	Tags      []Tag     `gorm:"many2many:champion_tags;" json:"tags,omitempty"`
	Versions  []Version `gorm:"many2many:champion_versions;" json:"versions,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Version struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
