package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDb() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to open env file %s", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		log.Printf("failed to open database %s", err)
	}

	fmt.Println("successfully connected to db")
	return db
}

