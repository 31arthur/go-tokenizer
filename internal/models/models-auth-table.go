package models

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/markbates/goth"
)

type UserAuth struct {
	ID                string `gorm:"primaryKey"`
	Provider          string `gorm:"not null"`
	UserID            string `gorm:"not null"`
	Username          string // Omitting constraints for nullable fields
	Email             string
	NameString        string
	AccessToken       string    `gorm:"not null"`
	RefreshToken      string    `gorm:"not null"`
	TokenExpiry       time.Time // GORM automatically maps time.Time to TIMESTAMP/TIMESTAMPTZ
	ProfilePictureURL string
	AuthTime          time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for this model
func (UserAuth) TableName() string {
	return "user_auth"
}

func MapUserToUserAuth(user goth.User) UserAuth {
	return UserAuth{
		ID:                GenerateHash(user.Provider, user.UserID),
		Provider:          user.Provider,
		UserID:            user.UserID,
		Username:          user.NickName, // Example: Mapping NickName to Username
		Email:             user.Email,
		NameString:        user.FirstName + " " + user.LastName,
		AccessToken:       user.AccessToken,
		RefreshToken:      user.RefreshToken,
		TokenExpiry:       user.ExpiresAt,
		ProfilePictureURL: user.AvatarURL,
		AuthTime:          time.Now(), // Set the AuthTime to the current time
	}
}

func GenerateHash(provider, userID string) string {
	hash := md5.New()
	hash.Write([]byte(provider + userID))
	return hex.EncodeToString(hash.Sum(nil))
}
