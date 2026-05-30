package model

import (
	"time"
)

type Registration struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name             string    `json:"name" gorm:"not null;size:50;uniqueIndex:idx_reg_unique,priority:1;index:idx_reg_name_class"`
	Gender           string    `json:"gender" gorm:"not null;size:2;check:gender IN ('男','女')"`
	BirthDate        time.Time `json:"birth_date" gorm:"not null;type:DATE"`
	Grade            string    `json:"grade" gorm:"not null;size:20"`
	ClassID          uint      `json:"class_id" gorm:"not null;uniqueIndex:idx_reg_unique,priority:3;index:idx_reg_class"`
	ParentName       string    `json:"parent_name" gorm:"not null;size:50"`
	ParentPhone      string    `json:"parent_phone" gorm:"not null;size:20;index"`
	Address          string    `json:"address" gorm:"not null;type:TEXT"`
	IDNumber         string    `json:"id_number" gorm:"not null;size:18;uniqueIndex:idx_reg_unique,priority:2;index"`
	PhotoPath        *string   `json:"photo_path" gorm:"size:255"`
	RegistrationTime time.Time `json:"registration_time" gorm:"autoCreateTime;index"`

	Class Class `json:"class" gorm:"foreignKey:ClassID"`
}

func (Registration) TableName() string {
	return "registrations"
}
