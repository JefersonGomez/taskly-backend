package main

import (
	"fmt"
	"os"
	"taskly-backend/controller"
	"taskly-backend/models"
	"taskly-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("No se pudo conectar a la base de datos...")
	}

	controller.DB = db

	models.MigrarTablas(db)

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":7000")

}
