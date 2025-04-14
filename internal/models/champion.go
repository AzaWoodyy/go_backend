package models

import (
	"time"

	"gorm.io/gorm"
)

type Champion struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	RiotID    string         `gorm:"type:varchar(255);uniqueIndex" json:"riot_id"`
	Key       string         `gorm:"type:varchar(255)" json:"key"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	Blurb     string         `gorm:"type:text" json:"blurb"`
	Tags      []Tag          `gorm:"many2many:champion_tags;" json:"tags,omitempty"`
	Versions  []Version      `gorm:"many2many:champion_versions;" json:"versions,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	gorm.Model
	Key       string    `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Version struct {
	gorm.Model
	Key       string    `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
