package database

import (
	"fmt"
	"log"
	"server/internal/models"
)

func (s *service) NewUserAuth() {
	s.tableCreator(&models.UserAuth{}, "User Auth")
	s.tableCreator(&models.UserDetails{}, "User Details")
}

func (s *service) tableCreator(tableModel any, tableName string) {
	migrator := s.db.Migrator()

	if !migrator.HasTable(tableModel) {
		err := s.db.AutoMigrate(tableModel)
		if err != nil {
			panic("failed to migrate table")
		}
		fmt.Println("Table " + tableName + " created successfully")
	} else {
		fmt.Println("Table " + tableName + " already exists")
	}
}

func (s *service) FreshCheck(uid string, email string) (bool, string, error) {
	var tempUAuthData models.UserAuth

	result := s.db.Where("user_id = ? OR email = ?", uid, email).First(&tempUAuthData)
	// fmt.Println("Result value", result, result.Error)
	if result.Error != nil {
		return false, "", result.Error
	}
	return true, tempUAuthData.ID, nil
}

func (s *service) InsertNewUser(user interface{}) error {

	if err := s.db.Create(user).Error; err != nil {
		log.Println("Error inserting new User Auth:", err)
		return err
	}
	return nil
}

func (s *service) SearchUserAuthByID(uid string) (models.UserAuth, error) {

	user, err := s.GetRowByID(uid)
	if err != nil {
		return models.UserAuth{}, err
	}

	return user, nil
}

func (s *service) ReplaceUserAuthByID(newUserAuth *models.UserAuth) error {
	var existingUserAuth models.UserAuth

	// Find the existing UserAuth record by ID
	result := s.db.Where("id = ?", newUserAuth.ID).First(&existingUserAuth)
	if result.Error != nil {
		return result.Error // Return error if record not found or any error occurs
	}

	// Update all fields with the new values
	existingUserAuth.LatestAuthTime = newUserAuth.LatestAuthTime
	existingUserAuth.RefreshToken = newUserAuth.RefreshToken
	existingUserAuth.AccessToken = newUserAuth.AccessToken
	existingUserAuth.RefreshTokenExpiry = newUserAuth.RefreshTokenExpiry
	// Update other fields as needed

	// Save the changes to the database
	if err := s.db.Save(&existingUserAuth).Error; err != nil {
		return err // Return error if update fails
	}

	return nil
}
