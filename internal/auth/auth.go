package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// create a new session

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	//connect the store to gothic store
	gothic.Store = store

	// you can get the provider names from the examples section of their github repository.
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/google/callback"),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), "http://localhost:3000/auth/github/callback"),
	)
}
