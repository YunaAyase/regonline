package repository

import (
	"regonline-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SiteSettingRepository struct {
	db *gorm.DB
}

func NewSiteSettingRepository(db *gorm.DB) *SiteSettingRepository {
	return &SiteSettingRepository{db: db}
}

func (r *SiteSettingRepository) GetAll() ([]model.SiteSetting, error) {
	var settings []model.SiteSetting
	if err := r.db.Order("id ASC").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *SiteSettingRepository) GetByKey(key string) (*model.SiteSetting, error) {
	var setting model.SiteSetting
	if err := r.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *SiteSettingRepository) Upsert(key, value string) error {
	setting := model.SiteSetting{
		Key:   key,
		Value: value,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&setting).Error
}

func (r *SiteSettingRepository) SeedDefaults() error {
	defaults := []struct {
		Key   string
		Value string
	}{
		{"site_name", "RegOnline 报名系统"},
		{"site_description", "在线报名管理系统"},
		{"icp_record", ""},
		{"copyright", "© 2025 RegOnline"},
	}

	for _, d := range defaults {
		var count int64
		r.db.Model(&model.SiteSetting{}).Where("key = ?", d.Key).Count(&count)
		if count == 0 {
			if err := r.db.Create(&model.SiteSetting{Key: d.Key, Value: d.Value}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}