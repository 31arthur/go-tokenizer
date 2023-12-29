package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"server/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService interface {
	Health() map[string]string
	InsertNewUser(*models.UserAuth) error
	UpdateUser(*models.UserAuth) error
}

type service struct {
	db *gorm.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
)

func NewCon() DBService {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return nil
	}
	s := &service{db: db}
	s.NewUserAuth()

	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var result int
	err := s.db.WithContext(ctx).Raw("SELECT 1").Row().Scan(&result)
	if err != nil {
		return map[string]string{
			"message": fmt.Sprintf("Database is down: %v", err),
		}
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
