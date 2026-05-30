package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"regonline-backend/internal/model"
	"regonline-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

type AuthService struct {
	adminRepo *repository.AdminRepository
}

func NewAuthService(adminRepo *repository.AdminRepository) *AuthService {
	return &AuthService{adminRepo: adminRepo}
}

func (s *AuthService) Login(username, password string) (*model.Admin, error) {
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return admin, nil
}

func (s *AuthService) SeedAdmin(username, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return s.adminRepo.SeedAdmin(username, string(bytes))
}

func (s *AuthService) UpdateAccount(currentUsername, newUsername, oldPassword, newPassword string) (*model.Admin, error) {
	admin, err := s.adminRepo.FindByUsername(currentUsername)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if newUsername != "" && newUsername != currentUsername {
		if _, err := s.adminRepo.FindByUsername(newUsername); err == nil {
			return nil, errors.New("新用户名已存在")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return nil, errors.New("原密码错误")
	}

	var hashedPassword string
	if newPassword != "" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		hashedPassword = string(bytes)
	}

	if err := s.adminRepo.Update(admin.ID, newUsername, hashedPassword); err != nil {
		return nil, fmt.Errorf("更新账号失败: %w", err)
	}

	admin.Username = newUsername
	admin.Password = hashedPassword
	return admin, nil
}

func (s *AuthService) MigratePasswords() error {
	admin, err := s.adminRepo.FindByUsername("admin")
	if err != nil {
		log.Println("No admin account found, skipping password migration")
		return nil
	}

	if strings.HasPrefix(admin.Password, "$2") {
		log.Println("Admin password is already bcrypt-hashed, skip migration")
		return nil
	}

	log.Println("Detected plain-text admin password, migrating to bcrypt...")
	plainPassword := admin.Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password during migration: %w", err)
	}

	if err := s.adminRepo.Update(admin.ID, "", string(bytes)); err != nil {
		return fmt.Errorf("failed to update password during migration: %w", err)
	}

	log.Println("Admin password migrated to bcrypt successfully")
	return nil
}