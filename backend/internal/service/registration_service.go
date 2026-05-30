package service

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"
	"regonline-backend/internal/util"
	errs "regonline-backend/internal/error"
)

type RegistrationService struct {
	regRepo    *repository.RegistrationRepository
	classRepo  *repository.ClassRepository
	photoDir   string
	writeMutex sync.Mutex
}

func NewRegistrationService(
	regRepo *repository.RegistrationRepository,
	classRepo *repository.ClassRepository,
	photoDir string,
) *RegistrationService {
	return &RegistrationService{
		regRepo:   regRepo,
		classRepo: classRepo,
		photoDir:  photoDir,
	}
}

type CreateRegistrationRequest struct {
	Name        string
	Gender      string
	BirthDate   time.Time
	Grade       string
	ClassID     uint
	ParentName  string
	ParentPhone string
	Address     string
	IDNumber    string
	Photo       []byte
	PhotoName   string
}

func (s *RegistrationService) Create(req *CreateRegistrationRequest) (*model.Registration, error) {
	req.IDNumber = util.NormalizeIDNumber(req.IDNumber)

	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()

	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	if err := s.validateIDNumber(req.IDNumber); err != nil {
		return nil, err
	}

	if err := s.validateGenderMatchesID(req.Gender, req.IDNumber); err != nil {
		return nil, err
	}

	if err := s.validateAge(req.BirthDate, req.ClassID); err != nil {
		return nil, err
	}

	if err := s.validateBirthDateMatchesID(req.BirthDate, req.IDNumber); err != nil {
		return nil, err
	}

	duplicate, err := s.regRepo.IsDuplicate(req.Name, req.IDNumber, req.ClassID)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	if duplicate {
		return nil, errs.NewDuplicateError("该学生已在此班级报名，请勿重复提交")
	}

	class, count, err := s.checkClassCapacity(req.ClassID)
	if err != nil {
		return nil, err
	}

	photoPath, err := s.savePhoto(req)
	if err != nil {
		return nil, fmt.Errorf("failed to save photo: %w", err)
	}

	registration := &model.Registration{
		Name:        req.Name,
		Gender:      req.Gender,
		BirthDate:   req.BirthDate,
		Grade:       req.Grade,
		ClassID:     req.ClassID,
		ParentName:  req.ParentName,
		ParentPhone: req.ParentPhone,
		Address:     req.Address,
		IDNumber:    req.IDNumber,
		PhotoPath:   photoPath,
	}

	if err := s.regRepo.Create(registration); err != nil {
		if photoPath != nil {
			os.Remove(*photoPath)
		}
		return nil, fmt.Errorf("failed to create registration: %w", err)
	}

	_ = count
	_ = class

	return registration, nil
}

func (s *RegistrationService) validateRequest(req *CreateRegistrationRequest) error {
	if req.Name == "" {
		return errs.NewValidationError("name", "学生姓名不能为空")
	}
	if req.Gender != "男" && req.Gender != "女" {
		return errs.NewValidationError("gender", "性别必须为男或女")
	}
	if req.BirthDate.IsZero() {
		return errs.NewValidationError("birth_date", "出生日期不能为空")
	}
	if req.Grade == "" {
		return errs.NewValidationError("grade", "年级不能为空")
	}
	if req.ClassID == 0 {
		return errs.NewValidationError("class_id", "请选择班级")
	}
	if req.ParentName == "" {
		return errs.NewValidationError("parent_name", "家长姓名不能为空")
	}
	if req.ParentPhone == "" {
		return errs.NewValidationError("parent_phone", "家长电话不能为空")
	}
	if req.Address == "" {
		return errs.NewValidationError("address", "家庭住址不能为空")
	}
	if req.IDNumber == "" {
		return errs.NewValidationError("id_number", "身份证号不能为空")
	}
	return nil
}

func (s *RegistrationService) validateIDNumber(idNumber string) error {
	if !util.ValidateIDNumberLength(idNumber) {
		return errs.NewValidationError("id_number", "身份证号必须为 18 位")
	}

	return nil
}

func (s *RegistrationService) validateGenderMatchesID(gender, idNumber string) error {
	if !util.ValidateGenderMatchesID(gender, idNumber) {
		genderDigit := idNumber[16] - '0'
		expectedGender := "女"
		if genderDigit%2 == 1 {
			expectedGender = "男"
		}
		return errs.NewValidationError("id_number",
			fmt.Sprintf("性别与身份证号不匹配！填写的性别：%s，身份证号第 17 位对应性别：%s",
				gender, expectedGender))
	}
	return nil
}

func (s *RegistrationService) validateAge(birthDate time.Time, classID uint) error {
	class, err := s.classRepo.FindByID(classID)
	if err != nil {
		return fmt.Errorf("failed to find class: %w", err)
	}

	age := util.CalculateAge(birthDate)

	if age < class.MinAge || age > class.MaxAge {
		return errs.NewValidationError("birth_date",
			fmt.Sprintf("年龄不符合要求！该班级要求 %d-%d 岁，当前年龄 %d 岁",
				class.MinAge, class.MaxAge, age))
	}

	return nil
}

func (s *RegistrationService) validateBirthDateMatchesID(birthDate time.Time, idNumber string) error {
	birthStr := birthDate.Format("20060102")
	if !util.ValidateBirthDateMatchesID(birthStr, idNumber) {
		idBirthStr := idNumber[6:14]
		return errs.NewValidationError("id_number",
			fmt.Sprintf("出生日期与身份证号不一致！填写的出生日期：%s-%s-%s，身份证号中的出生日期：%s-%s-%s",
				birthStr[0:4], birthStr[4:6], birthStr[6:8],
				idBirthStr[0:4], idBirthStr[4:6], idBirthStr[6:8]))
	}
	return nil
}

func (s *RegistrationService) checkClassCapacity(classID uint) (*model.Class, int64, error) {
	class, err := s.classRepo.FindByID(classID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find class: %w", err)
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

func (s *RegistrationService) savePhoto(req *CreateRegistrationRequest) (*string, error) {
	if req.Photo == nil || len(req.Photo) == 0 {
		return nil, nil
	}

	if err := os.MkdirAll(s.photoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create photo directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405_000000")
	hash := fmt.Sprintf("%x", md5.Sum([]byte(req.Name)))[:8]
	ext := filepath.Ext(req.PhotoName)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%s%s", timestamp, hash, ext)
	filePath := filepath.Join(s.photoDir, filename)

	if err := os.WriteFile(filePath, req.Photo, 0644); err != nil {
		return nil, fmt.Errorf("failed to write photo file: %w", err)
	}

	return &filePath, nil
}

func (s *RegistrationService) ListRegistrations() ([]model.Registration, error) {
	return s.regRepo.FindAll()
}

func (s *RegistrationService) GetRegistrationByID(id uint) (*model.Registration, error) {
	reg, err := s.regRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find registration: %w", err)
	}
	return reg, nil
}

func (s *RegistrationService) SearchRegistrations(name string, classID uint) ([]model.Registration, error) {
	return s.regRepo.Search(name, classID)
}

func (s *RegistrationService) DeleteRegistration(id uint) error {
	reg, err := s.regRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find registration: %w", err)
	}

	if reg.PhotoPath != nil {
		os.Remove(*reg.PhotoPath)
	}

	if err := s.regRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}

	return nil
}

func (s *RegistrationService) GetStats() (map[string]int64, error) {
	classes, err := s.classRepo.FindAll()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, class := range classes {
		count, _ := s.classRepo.CountRegistrations(class.ID)
		stats[class.Name] = count
	}

	return stats, nil
}

func (s *RegistrationService) GetDateRangeForClass(classID uint) (time.Time, time.Time, error) {
	class, err := s.classRepo.FindByID(classID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	minDate, maxDate := util.CalculateDateRangeByAge(class.MinAge, class.MaxAge)
	return minDate, maxDate, nil
}

func (s *RegistrationService) GetDefaultDateRange() (time.Time, time.Time) {
	return util.CalculateDateRangeByAge(5, 18)
}