package middleware

import (
	"fmt"
	"net/http"
	"os"
	"server/internal/database"
	"server/internal/tokenizer"
	"server/internal/utility"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthCheck(db database.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err2 := c.Cookie("pmt_auth_acct")
		fmt.Println("token Value @ ", tokenString)
		if err2 != nil {
			// fmt.Println("cookie Stringization Error", err2)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Validate the token and if valid get the uid, exp
		uid, exp, err := utility.ValidateToken(tokenString, os.Getenv("HMAC_SECT"))
		fmt.Println("Error while validating token:", err)
		if err != nil {
			fmt.Println("Uid and Exp")
			if float64(time.Now().Unix()) > exp {
				fmt.Println("Ref Token Uid", uid)
				accessToken, errT := RefTokenValidate(db, uid)
				if errT != nil {
					fmt.Println("Ref Token Fetch Error @ ", errT)
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				} else {
					c.SetCookie("pmt_auth_acct", accessToken, 3600, "/", "localhost", false, true)
					finalStep(db, c, uid)
				}

			}
			fmt.Println("token Validation Error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		finalStep(db, c, uid)
	}
}

func finalStep(db database.DBService, c *gin.Context, uid string) {
	user, err1 := db.SearchUserAuthByID(uid)
	if err1 != nil {
		// fmt.Println("Db Error @", err1)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set user in the context for other handlers to use
	c.Set("user", &user)

	c.Next()
}

func RefTokenValidate(db database.DBService, uid string) (string, error) {
	refToken, _, err := db.GetRefreshTokenDetailsByID(uid)
	if err != nil {
		return "", err
	}

	_, _, errV := utility.ValidateToken(refToken, os.Getenv("HMAC_RF_SECT"))
	fmt.Println("Error while validating token:", err)
	if errV != nil {
		return "", err
	}

	accessToken, errAT := tokenizer.AccessTokenGenerator(uid)
	if errAT != nil {
		return "", err
	}

	errDB := db.PutNewAccessTokenByID(uid, accessToken)
	if errDB != nil {
		return "", errDB
	}

	return accessToken, nil
}
