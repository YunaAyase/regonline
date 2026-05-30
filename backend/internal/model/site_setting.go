package model

import (
	"time"
)

type SiteSetting struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Key       string    `gorm:"uniqueIndex;not null;size:100"`
	Value     string    `gorm:"not null;size:4096"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}