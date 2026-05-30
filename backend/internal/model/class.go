package model

import "time"

type Class struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Description string    `json:"description" gorm:"not null;default:'';size:255"`
	MaxStudents int       `json:"max_students" gorm:"not null;default:30"`
	MinAge      int       `json:"min_age" gorm:"not null;default:5"`
	MaxAge      int       `json:"max_age" gorm:"not null;default:18"`
	Enabled     bool      `json:"enabled" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Class) TableName() string {
	return "classes"
}