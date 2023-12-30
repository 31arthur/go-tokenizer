package database

import (
	"time"
)

func (s *service) GetAccessTokenByID(id string) (string, error) {
	user, err := s.GetRowByID(id)

	if err != nil {
		return "", err
	}

	accessToken := user.AccessToken
	return accessToken, nil
}

func (s *service) GetRefreshTokenDetailsByID(id string) (string, time.Time, error) {
	user, err := s.GetRowByID(id)
	if err != nil {
		return "", time.Time{}, err
	}
	return user.RefreshToken, user.RefreshTokenExpiry, nil
}

func (s *service) PutNewAccessTokenByID(uid string, accessToken string) error {

	user, errR := s.GetRowByID(uid)

	if errR != nil {
		return errR
	}
	user.AccessToken = accessToken

	// Save the changes to the database
	if err := s.db.Save(&user).Error; err != nil {
		return err // Return error if update fails
	}

	return nil
}
