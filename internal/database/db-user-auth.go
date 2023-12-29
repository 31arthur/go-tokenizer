package database

import (
	"log"
	"server/internal/models"
	"time"
)

func (s *service) NewUserAuth() {
	migrator := s.db.Migrator()

	if !migrator.HasTable(&models.UserAuth{}) {
		err := s.db.AutoMigrate(&models.UserAuth{})
		if err != nil {
			panic("failed to migrate table")
		}
		println("Table UserAuth created successfully")
	} else {
		println("Table UserAuth already exists")
	}
}

func (s *service) InsertNewUser(user *models.UserAuth) error {
	var existingUser models.UserAuth

	// Check if User ID or Email ID already exists in the table
	result := s.db.Where("user_id = ? OR email = ?", user.UserID, user.Email).First(&existingUser)
	if result.Error == nil {
		// User exists; update the user
		return s.UpdateUser(user)
	}

	// User does not exist; proceed with insertion
	if err := s.db.Create(user).Error; err != nil {
		log.Println("Error inserting new user:", err)
		return err
	}
	return nil
}

func (s *service) UpdateUser(user *models.UserAuth) error {
	// Fetch the existing user record
	var existingUser models.UserAuth
	if err := s.db.Where("user_id = ? OR email = ?", user.UserID, user.Email).First(&existingUser).Error; err != nil {
		log.Println(err)
		return err
	}

	updateFields := make(map[string]interface{})
	updateFields["provider"] = user.Provider
	updateFields["access_token"] = user.AccessToken
	updateFields["refresh_token"] = user.RefreshToken
	updateFields["token_expiry"] = user.TokenExpiry
	updateFields["auth_time"] = time.Now()

	// Update name if it's not blank
	if user.NameString != "" {
		updateFields["name_string"] = user.NameString
	}

	if err := s.db.Model(&existingUser).Updates(updateFields).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}
