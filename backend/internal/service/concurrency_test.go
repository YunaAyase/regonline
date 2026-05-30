package service

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegistrationService_CapacityConcurrent(t *testing.T) {
	svc, db := setupTestService(t)
	classID := seedTestClass(t, db, 3, 7, 16)

	birthDate := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)

	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	names := []string{"测试学生A", "测试学生B", "测试学生C", "测试学生D", "测试学生E"}
	idNumbers := []string{"411728201503150012", "411728201503150020", "411728201503150039", "411728201503150047", "411728201503150055"}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			req := &CreateRegistrationRequest{
				Name:        names[idx],
				Gender:      "男",
				BirthDate:   birthDate,
				Grade:       "三年级",
				ClassID:     classID,
				ParentName:  "张某某",
				ParentPhone: "13800138000",
				Address:     "XX小区1栋",
				IDNumber:    idNumbers[idx],
			}

			_, err := svc.Create(req)
			mu.Lock()
			if err == nil {
				successCount++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 3, successCount, "班级容量为 3，应该只有 3 人报名成功")

	sqlDB, _ := db.DB()
	var count int64
	sqlDB.QueryRow("SELECT count(*) FROM registrations").Scan(&count)
	assert.Equal(t, int64(3), count, "数据库中应该只有 3 条记录")
}
