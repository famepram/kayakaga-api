package main

import (
	"kayakaga-api/di"
	"kayakaga-api/domain/mysql"
	"kayakaga-api/utils/router"
	"log"
	"os"

	_ "kayakaga-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Kayakaga API
// @version 1.0
// @description Personal Finance Advisor REST API built with Go, Gin, GORM, and MySQL
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Auth
// @tag.description Authentication endpoints

// @tag.name Users
// @tag.description User profile management

// @tag.name Accounts
// @tag.description Account management

// @tag.name Transactions
// @tag.description Transaction management

// @tag.name Goals
// @tag.description Goal tracking

// @tag.name Analytics
// @tag.description Financial analytics

// @tag.name Simulation
// @tag.description Investment simulation

// @tag.name Master Data
// @tag.description Master data endpoints

// @tag.name Chat
// @tag.description AI Chat endpoints

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := mysql.NewConfig()
	db, err := mysql.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	container := di.InitializeContainer(db)
	r := router.SetupRouter(container)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Kayakaga API on port %s (%s mode)", port, appEnv)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
