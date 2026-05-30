package model

import (
	"time"
)

type Admin struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"uniqueIndex;not null;size:50"`
	Password  string    `gorm:"not null;size:128"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Admin) TableName() string {
	return "admins"
}
