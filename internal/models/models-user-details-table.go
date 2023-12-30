package models

import (
	"github.com/markbates/goth"
)

type UserDetails struct {
	ID                string `gorm:"primaryKey"`
	Provider          string `gorm:"not null"`
	UserID            string `gorm:"not null"`
	Username          string // Omitting constraints for nullable fields
	Email             string
	FirstName         string
	LastName          string
	ProfilePictureURL string
	Phone             string
	GitHubAdd         string
	Description       string

	UserAuth UserAuth `gorm:"foreignKey:ID"`
}

func (UserDetails) TableName() string {
	return "users_details"
}

func MapUserToUserDetails(user goth.User, uid string) UserDetails {
	return UserDetails{
		ID:                uid,
		Provider:          user.Provider,
		UserID:            user.UserID,
		Username:          user.NickName, // Example: Mapping NickName to Username
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		ProfilePictureURL: user.AvatarURL,
	}
}
