package utility

import (
	"fmt"
	"server/internal/tokenizer"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AllTokens(id string) (string, string, error) {

	accessToken, errAT := tokenizer.AccessTokenGenerator(id)
	refreshToken, errRT := tokenizer.RefreshTokenGenerator(id)

	if errAT != nil {
		return "", "", errAT
	}

	if errRT != nil {
		return "", "", errRT
	}

	return accessToken, refreshToken, nil

}

func GetTokenExpiration(tokenString string) (time.Time, error) {

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse token: %s", err)
	}

	// Check if the token claims can be asserted as MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid token claims")
	}

	// Retrieve the expiration time from the token claims
	exp, ok := claims["exp"].(float64)
	if !ok {
		return time.Time{}, fmt.Errorf("expiration time not found or invalid")
	}

	// Convert the expiration time from Unix timestamp to time.Time
	expirationTime := time.Unix(int64(exp), 0)

	return expirationTime, nil
}

func ValidateToken(tokenString string, secret string) (string, float64, error) {
	// Parse the token string using the specified parsing function
	fmt.Println("Validate Token Called")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the subject and expiration claims
		subject := fmt.Sprintf("%v", claims["sub"]) // Assuming "sub" is a string claim
		expiration := claims["exp"].(float64)       // Assuming "exp" is a float64 claim

		fmt.Println("Subject:", subject)
		fmt.Println("Expiration:", time.Unix(int64(expiration), 0).Format(time.RFC3339))

		return subject, expiration, nil

	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Token is expired but claims can still be accessed
		subject := fmt.Sprintf("%v", claims["sub"])
		expiration := claims["exp"].(float64)

		fmt.Println("Subject (Expired):", subject)
		fmt.Println("Expiration (Expired):", time.Unix(int64(expiration), 0).Format(time.RFC3339))

		return subject, expiration, err

	} else {
		// Token is not valid
		fmt.Println("Invalid token")
		return "", 0.0, err
	}
}
