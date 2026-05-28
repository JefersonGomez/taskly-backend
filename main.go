// @title           Taskly API
// @version         1.0.0
// @description     Documentación automática del backend de Taskly
// @host            localhost:7000
// @BasePath        /
// @schemes         http

package main

import (
	"fmt"
	"os"
	"taskly-backend/controller"
	_ "taskly-backend/docs"
	"taskly-backend/models"
	"taskly-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load()

	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	fmt.Println("DSN usado:", dsn) // log temporal

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos...")
	}

	controller.DB = db
	models.MigrarTablas(db)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://taskly-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	r.Run(":" + port)
}
