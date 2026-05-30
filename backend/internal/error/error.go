package errs

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("字段 %s: %s", e.Field, e.Message)
	}
	return e.Message
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

type NotFoundError struct {
	Resource string
	ID       any
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s (ID: %v) 不存在", e.Resource, e.ID)
}

func NewNotFoundError(resource string, id any) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

type DuplicateError struct {
	Message string
}

func (e *DuplicateError) Error() string {
	return e.Message
}

func NewDuplicateError(message string) *DuplicateError {
	return &DuplicateError{
		Message: message,
	}
}

type CapacityError struct {
	ClassName string
	Current   int
	Max       int
}

func (e *CapacityError) Error() string {
	return fmt.Sprintf("%s 已报满（%d/%d）", e.ClassName, e.Current, e.Max)
}

func NewCapacityError(className string, current, max int) *CapacityError {
	return &CapacityError{
		ClassName: className,
		Current:   current,
		Max:       max,
	}
}

var (
	ErrClassNotFound     = errors.New("班级不存在")
	ErrRegistrationExists = errors.New("该学生已在此班级报名")
)

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
