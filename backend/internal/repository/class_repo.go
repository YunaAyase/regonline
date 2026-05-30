package repository

import (
	"regonline-backend/internal/model"

	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

func (r *ClassRepository) FindAll() ([]model.Class, error) {
	var classes []model.Class
	err := r.db.Order("id ASC").Find(&classes).Error
	return classes, err
}

func (r *ClassRepository) FindByID(id uint) (*model.Class, error) {
	var class model.Class
	err := r.db.First(&class, id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *ClassRepository) FindByName(name string) (*model.Class, error) {
	var class model.Class
	err := r.db.Where("name = ?", name).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *ClassRepository) Create(class *model.Class) error {
	return r.db.Create(class).Error
}

func (r *ClassRepository) Delete(id uint) error {
	return r.db.Delete(&model.Class{}, id).Error
}

func (r *ClassRepository) CountRegistrations(classID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Registration{}).
		Where("class_id = ?", classID).
		Count(&count).Error
	return count, err
}

func (r *ClassRepository) Update(class *model.Class) error {
	return r.db.Save(class).Error
}
