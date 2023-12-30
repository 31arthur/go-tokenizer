package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"server/internal/utility"
	"time"

	"github.com/markbates/goth"
)

type UserAuth struct {
	ID                 string `gorm:"primaryKey"`
	Provider           string `gorm:"not null"`
	UserID             string `gorm:"not null"`
	Email              string
	AccessToken        string    `gorm:"not null"`
	RefreshToken       string    `gorm:"not null"`
	RefreshTokenExpiry time.Time // GORM automatically maps time.Time to TIMESTAMP/TIMESTAMPTZ
	LatestAuthTime     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for this model
func (UserAuth) TableName() string {
	return "user_auth"
}

func MapUserToUserAuth(user goth.User, uid string) (UserAuth, error) {

	id := uid
	if uid == "" {
		id = GenerateHash(user.Provider, user.UserID)
	}

	accToken, refToken, err := utility.AllTokens(id)
	rtExpiry, err1 := utility.GetTokenExpiration(refToken)

	if err != nil || err1 != nil {
		return UserAuth{}, fmt.Errorf("error occurred: err=%v, err1=%v", err, err1)
	}

	newUserAuth := UserAuth{
		ID:                 id,
		Provider:           user.Provider,
		UserID:             user.UserID,
		Email:              user.Email,
		AccessToken:        accToken,
		RefreshToken:       refToken,
		RefreshTokenExpiry: rtExpiry,
		LatestAuthTime:     time.Now(), // Set the AuthTime to the current time
	}

	return newUserAuth, nil
}

func GenerateHash(provider, userID string) string {
	hash := md5.New()
	hash.Write([]byte(provider + userID))
	return hex.EncodeToString(hash.Sum(nil))
}
