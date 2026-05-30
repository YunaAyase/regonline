package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"regonline-backend/internal/config"
	"regonline-backend/internal/database"
	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	tmpFile := filepath.Join(t.TempDir(), "test.db")

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:      "sqlite",
			Path:        tmpFile,
			AutoMigrate: true,
		},
	}

	db, err := database.NewDB(cfg)
	require.NoError(t, err)

	err = database.AutoMigrate(db, cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		// Give Windows time to release the file lock
		time.Sleep(100 * time.Millisecond)
		os.Remove(tmpFile)
	})

	return db
}

func setupTestService(t *testing.T) (*RegistrationService, *gorm.DB) {
	db := setupTestDB(t)

	classRepo := repository.NewClassRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	photoDir := t.TempDir()

	return NewRegistrationService(regRepo, classRepo, photoDir), db
}

func seedTestClass(t *testing.T, db *gorm.DB, maxStudents, minAge, maxAge int) uint {
	class := model.Class{
		Name:        "测试班",
		MaxStudents: maxStudents,
		MinAge:      minAge,
		MaxAge:      maxAge,
	}
	db.Create(&class)
	return class.ID
}

func TestRegistrationService_Create_Success(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728201503150012",
	}

	reg, err := svc.Create(req)
	require.NoError(t, err)
	assert.Equal(t, "测试学生", reg.Name)
	assert.Equal(t, classID, reg.ClassID)
	assert.Equal(t, "411728201503150012", reg.IDNumber)
}

func TestRegistrationService_Create_Duplicate(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728201503150012",
	}

	_, err := svc.Create(req)
	require.NoError(t, err)

	_, err = svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "该学生已在此班级报名")
}

func TestRegistrationService_Create_AgeTooYoung(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生B",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "一年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728202003150010",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "年龄不符合要求")
}

func TestRegistrationService_Create_AgeTooOld(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2005, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生C",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "初三",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728200503150018",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "年龄不符合要求")
}

func TestRegistrationService_Create_BirthDateMismatch(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生D",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728202001010014",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "出生日期与身份证号不一致")
}

func TestRegistrationService_Create_InvalidIDNumberLength(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生E",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "4117282015031",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "身份证号必须为 18 位")
}

func TestRegistrationService_Create_InvalidIDNumberFormat(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生F",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "41172820150315001a",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "身份证号格式不正确")
}

func TestRegistrationService_Create_InvalidIDNumberChecksum(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 30, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	req := &CreateRegistrationRequest{
		Name:        "测试学生G",
		Gender:      "男",
		BirthDate:   birthDate,
		Grade:       "三年级",
		ClassID:     classID,
		ParentName:  "张某某",
		ParentPhone: "13800138000",
		Address:     "XX小区1栋",
		IDNumber:    "411728201503150013",
	}

	_, err := svc.Create(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "身份证号校验码不正确")
}
