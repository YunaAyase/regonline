package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassService_ListClasses(t *testing.T) {
	svc, db := setupTestServiceForClass(t)

	seedClass(t, db, "编程班", 30, 7, 16)
	seedClass(t, db, "手工班", 25, 6, 15)

	classes, err := svc.ListClasses()
	require.NoError(t, err)
	assert.Len(t, classes, 2)
	assert.Equal(t, "编程班", classes[0].Name)
	assert.Equal(t, int64(0), classes[0].CurrentCount)
}

func TestClassService_GetClassByName(t *testing.T) {
	svc, db := setupTestServiceForClass(t)
	seedClass(t, db, "编程班", 30, 7, 16)

	classInfo, err := svc.GetClassByName("编程班")
	require.NoError(t, err)
	assert.Equal(t, "编程班", classInfo.Name)
	assert.Equal(t, 30, classInfo.MaxStudents)
	assert.Equal(t, 7, classInfo.MinAge)
	assert.Equal(t, 16, classInfo.MaxAge)
}

func TestClassService_GetClassByName_NotFound(t *testing.T) {
	svc, _ := setupTestServiceForClass(t)

	_, err := svc.GetClassByName("不存在的班级")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}
