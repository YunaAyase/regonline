package service

import (
	"regonline-backend/internal/repository"
)

type SiteSettingService struct {
	repo *repository.SiteSettingRepository
}

func NewSiteSettingService(repo *repository.SiteSettingRepository) *SiteSettingService {
	return &SiteSettingService{repo: repo}
}

func (s *SiteSettingService) GetSettings() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(settings))
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}

func (s *SiteSettingService) UpdateSettings(updates map[string]string) error {
	for key, value := range updates {
		if err := s.repo.Upsert(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *SiteSettingService) SeedDefaults() error {
	return s.repo.SeedDefaults()
}