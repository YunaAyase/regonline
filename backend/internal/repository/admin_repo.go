package repository

import (
	"regonline-backend/internal/model"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) AdminCount() int64 {
	var count int64
	r.db.Model(&model.Admin{}).Count(&count)
	return count
}

func (r *AdminRepository) SeedAdmin(username, password string) error {
	var count int64
	r.db.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	admin := model.Admin{
		Username: username,
		Password: password,
	}
	return r.db.Create(&admin).Error
}

func (r *AdminRepository) Update(id uint, username, password string) error {
	updates := map[string]interface{}{}
	if username != "" {
		updates["username"] = username
	}
	if password != "" {
		updates["password"] = password
	}
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Updates(updates).Error
}
