package tokenizer

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthToken(val string, valTime int, secret string) (string, error) {

	secretKey := []byte(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": val,
		"exp": time.Now().Add(time.Hour * time.Duration(valTime)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	//secretKey should converted into byte array. Remember that.
	tokenString, err := token.SignedString(secretKey)

	// fmt.Println("token string with error", tokenString, err, secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func AccessTokenGenerator(uid string) (string, error) {
	accessToken, err := AuthToken(uid, 2, os.Getenv("HMAC_SECT"))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func RefreshTokenGenerator(uid string) (string, error) {
	refreshToken, err := AuthToken(uid, 24, os.Getenv("HMAC_RF_SECT"))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
