package services

import (
	"errors"

	"github.com/fitrianabila2025group/videoxnx/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// EnsureAdminUser creates the default admin if not present, or updates the
// stored password if it no longer matches the configured ADMIN_PASSWORD.
// This makes credentials in .env the single source of truth for production.
func EnsureAdminUser(db *gorm.DB, email, password string) error {
	if email == "" || password == "" {
		return errors.New("admin email/password required")
	}
	var existing models.User
	err := db.Where("email = ?", email).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u := models.User{Email: email, PasswordHash: string(hash), Role: "admin"}
		return db.Create(&u).Error
	}
	// User exists: refresh password if it changed.
	if bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(password)) != nil {
		return db.Model(&existing).Updates(map[string]any{
			"password_hash": string(hash),
			"role":          "admin",
		}).Error
	}
	return nil
}

func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
