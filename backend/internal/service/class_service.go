package service

import (
	"errors"

	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"
	errs "regonline-backend/internal/error"

	"gorm.io/gorm"
)

type ClassService struct {
	classRepo    *repository.ClassRepository
	regRepo      *repository.RegistrationRepository
}

func NewClassService(classRepo *repository.ClassRepository, regRepo *repository.RegistrationRepository) *ClassService {
	return &ClassService{
		classRepo: classRepo,
		regRepo:   regRepo,
	}
}

type ClassInfo struct {
	model.Class
	CurrentCount int64 `json:"current_count"`
}

func (s *ClassService) ListClasses() ([]ClassInfo, error) {
	classes, err := s.classRepo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]ClassInfo, 0, len(classes))
	for _, class := range classes {
		count, _ := s.classRepo.CountRegistrations(class.ID)
		result = append(result, ClassInfo{
			Class:        class,
			CurrentCount: count,
		})
	}

	return result, nil
}

func (s *ClassService) GetClassByID(id uint) (*ClassInfo, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("班级", id)
		}
		return nil, err
	}

	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) GetClassByName(name string) (*ClassInfo, error) {
	class, err := s.classRepo.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("班级", name)
		}
		return nil, err
	}

	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) UpdateClass(id uint, maxStudents, minAge, maxAge int) (*ClassInfo, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("班级", id)
		}
		return nil, err
	}

	class.MaxStudents = maxStudents
	class.MinAge = minAge
	class.MaxAge = maxAge

	if err := s.classRepo.Update(class); err != nil {
		return nil, err
	}

	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) CreateClass(name, description string, maxStudents, minAge, maxAge int) (*ClassInfo, error) {
	class := &model.Class{
		Name:        name,
		Description: description,
		MaxStudents: maxStudents,
		MinAge:      minAge,
		MaxAge:      maxAge,
		Enabled:     true,
	}
	if err := s.classRepo.Create(class); err != nil {
		return nil, errs.NewDuplicateError("班级名称已存在")
	}
	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) UpdateClassFull(id uint, name, description string, maxStudents, minAge, maxAge int) (*ClassInfo, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("班级", id)
		}
		return nil, err
	}

	class.Name = name
	class.Description = description
	class.MaxStudents = maxStudents
	class.MinAge = minAge
	class.MaxAge = maxAge

	if err := s.classRepo.Update(class); err != nil {
		return nil, errs.NewDuplicateError("班级名称已存在")
	}

	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) ToggleClass(id uint, enabled bool) (*ClassInfo, error) {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFoundError("班级", id)
		}
		return nil, err
	}

	class.Enabled = enabled
	if err := s.classRepo.Update(class); err != nil {
		return nil, err
	}

	count, _ := s.classRepo.CountRegistrations(class.ID)
	return &ClassInfo{
		Class:        *class,
		CurrentCount: count,
	}, nil
}

func (s *ClassService) DeleteClass(id uint) error {
	count, err := s.classRepo.CountRegistrations(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errs.NewValidationError("class_id", "该班级还有报名记录，无法删除")
	}
	return s.classRepo.Delete(id)
}

func (s *ClassService) CheckCapacity(classID uint) (*model.Class, int64, error) {
	class, err := s.classRepo.FindByID(classID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errs.NewNotFoundError("班级", classID)
		}
		return nil, 0, err
	}

	count, err := s.classRepo.CountRegistrations(classID)
	if err != nil {
		return nil, 0, err
	}

	if count >= int64(class.MaxStudents) {
		return nil, count, errs.NewCapacityError(class.Name, int(count), class.MaxStudents)
	}

	return class, count, nil
}
