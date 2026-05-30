package repository

import (
	"strings"

	"regonline-backend/internal/model"

	"gorm.io/gorm"
)

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) FindAll() ([]model.Registration, error) {
	var registrations []model.Registration
	err := r.db.Preload("Class").
		Order("registration_time DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) FindByID(id uint) (*model.Registration, error) {
	var reg model.Registration
	err := r.db.Preload("Class").First(&reg, id).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *RegistrationRepository) FindByName(name string) ([]model.Registration, error) {
	var registrations []model.Registration
	err := r.db.Preload("Class").
		Where("name LIKE ?", "%"+name+"%").
		Order("registration_time DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) FindByClassID(classID uint) ([]model.Registration, error) {
	var registrations []model.Registration
	err := r.db.Preload("Class").
		Where("class_id = ?", classID).
		Order("registration_time DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) FindByNameAndClass(name string, classID uint) ([]model.Registration, error) {
	var registrations []model.Registration
	err := r.db.Preload("Class").
		Where("name = ? AND class_id = ?", name, classID).
		Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) FindByIDNumberAndClass(idNumber string, classID uint) ([]model.Registration, error) {
	var registrations []model.Registration
	err := r.db.Preload("Class").
		Where("id_number = ? AND class_id = ?", idNumber, classID).
		Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) Search(name string, classID uint) ([]model.Registration, error) {
	query := r.db.Model(&model.Registration{}).Preload("Class")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if classID > 0 {
		query = query.Where("class_id = ?", classID)
	}

	var registrations []model.Registration
	err := query.Order("registration_time DESC").Find(&registrations).Error
	return registrations, err
}

func (r *RegistrationRepository) Create(reg *model.Registration) error {
	return r.db.Create(reg).Error
}

func (r *RegistrationRepository) IsDuplicate(name, idNumber string, classID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Registration{}).
		Where(
			"(name = ? OR id_number = ?) AND class_id = ?",
			name, idNumber, classID,
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RegistrationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Registration{}, id).Error
}

func (r *RegistrationRepository) DeleteByFile(photoPath string) error {
	if photoPath == "" {
		return nil
	}

	var regs []model.Registration
	if err := r.db.Where("photo_path = ?", photoPath).Find(&regs).Error; err != nil {
		return err
	}

	for _, reg := range regs {
		reg.PhotoPath = nil
		if err := r.db.Save(&reg).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RegistrationRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Registration{}).Count(&count).Error
	return count, err
}

func (r *RegistrationRepository) CountByClassID(classID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Registration{}).
		Where("class_id = ?", classID).
		Count(&count).Error
	return count, err
}

func (r *RegistrationRepository) FindDuplicates() ([]model.Registration, error) {
	var registrations []model.Registration

	var subQuery string
	switch {
	case strings.Contains(r.db.Dialector.Name(), "sqlite"):
		subQuery = `SELECT name, id_number, class_id FROM registrations GROUP BY name, id_number, class_id HAVING COUNT(*) > 1`
	default:
		subQuery = `SELECT name, id_number, class_id FROM registrations GROUP BY name, id_number, class_id HAVING COUNT(*) > 1`
	}

	err := r.db.Preload("Class").
		Where("(name, id_number, class_id) IN (?)", r.db.Raw(subQuery)).
		Order("name, class_id, registration_time").
		Find(&registrations).Error
	return registrations, err
}
