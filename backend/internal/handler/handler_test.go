package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"regonline-backend/internal/config"
	"regonline-backend/internal/database"
	"regonline-backend/internal/handler"
	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"
	"regonline-backend/internal/router"
	"regonline-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedTestDBClass(t *testing.T, db *gorm.DB, name string, maxStudents, minAge, maxAge int) uint {
	class := model.Class{
		Name:        name,
		MaxStudents: maxStudents,
		MinAge:      minAge,
		MaxAge:      maxAge,
	}
	db.Create(&class)
	return class.ID
}

func setupTestServer(t *testing.T) (*httptest.Server, *gorm.DB) {
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

	classRepo := repository.NewClassRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	photoDir := t.TempDir()

	classService := service.NewClassService(classRepo, regRepo)
	regService := service.NewRegistrationService(regRepo, classRepo, photoDir)
	ocrService := service.NewOCRService()

	classHandler := handler.NewClassHandler(classService)
	regHandler := handler.NewRegistrationHandler(regService)
	ocrHandler := handler.NewOCRHandler(ocrService)
	statsHandler := handler.NewStatsHandler(regService)

	r := router.SetupRouter(classHandler, regHandler, ocrHandler, statsHandler)

	server := httptest.NewServer(r)
	t.Cleanup(func() {
		server.Close()
		sqlDB, _ := db.DB()
		sqlDB.Close()
		time.Sleep(100 * time.Millisecond)
		os.Remove(tmpFile)
	})

	return server, db
}

func TestHandler_HealthCheck(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestHandler_ListClasses(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	seedTestDBClass(t, db, "编程班", 30, 7, 16)
	seedTestDBClass(t, db, "手工班", 25, 6, 15)

	resp, err := http.Get(server.URL + "/api/classes")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Code int `json:"code"`
		Data []struct {
			Name         string `json:"name"`
			CurrentCount int64  `json:"current_count"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)
	assert.Len(t, result.Data, 2)
}

func TestHandler_CreateRegistration_Success(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	classID := seedTestDBClass(t, db, "编程班", 30, 7, 16)

	body := `name=测试学生&gender=男&birth_date=2015-03-15&grade=三年级&class_id=` + fmt.Sprintf("%d", classID) + `&parent_name=张某某&parent_phone=13800138000&address=XX小区1栋&id_number=411728201503150012`

	resp, err := http.Post(server.URL+"/api/registrations", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 0, result.Code)
	assert.Equal(t, "created", result.Message)
}

func TestHandler_CreateRegistration_Duplicate(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	classID := seedTestDBClass(t, db, "编程班", 30, 7, 16)

	body := `name=测试学生&gender=男&birth_date=2015-03-15&grade=三年级&class_id=` + fmt.Sprintf("%d", classID) + `&parent_name=张某某&parent_phone=13800138000&address=XX小区1栋&id_number=411728201503150012`

	resp, _ := http.Post(server.URL+"/api/registrations", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	resp.Body.Close()

	resp, _ = http.Post(server.URL+"/api/registrations", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	assert.Equal(t, 409, result.Code)
	assert.Contains(t, result.Message, "该学生已在此班级报名")
}

func TestHandler_CreateRegistration_AgeTooYoung(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	classID := seedTestDBClass(t, db, "编程班", 30, 7, 16)

	body := `name=测试学生B&gender=男&birth_date=2020-03-15&grade=一年级&class_id=` + fmt.Sprintf("%d", classID) + `&parent_name=张某某&parent_phone=13800138000&address=XX小区1栋&id_number=411728202003150010`

	resp, err := http.Post(server.URL+"/api/registrations", "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, 400, result.Code)
	assert.Contains(t, result.Message, "年龄不符合要求")
}
