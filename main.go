package main

import (
	"fmt"
	"go-gin-postgre/domain"
	"go-gin-postgre/domain/usecases"
	"go-gin-postgre/infrastructures"
	"go-gin-postgre/interfaces/handlers"
	"go-gin-postgre/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local" // 기본값
	}

	// 해당 환경의 .env 파일 불러오기
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatalf("Error loading .env.%s file", env)
	}

	middleware.InitJWTSecretKey()
	handlers.InitGoogleOauthConfig()
	initDB()

	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		log.Printf("googleRedirectURL: %v", googleRedirectURL)
	}

	// 의존성 주입
	repo := infrastructures.NewUserRepository(db)
	usecase := usecases.NewUserUsecase(repo)
	handler := handlers.NewUserHandler(usecase)

	r := gin.Default()

	r.GET("/users", middleware.JwtAuthMiddleware(), handler.GetUsers)
	r.GET("/users/:id", handler.GetUser)
	r.POST("/users", handler.CreateUser)
	r.PUT("/users/:id", handler.UpdateUser)
	r.DELETE("/users/:id", handler.DeleteUser)

	r.GET("/auth/google/login", handlers.GoogleLoginHandler)
	r.GET("/auth/google/callback", handlers.GoogleCallbackHandler)

	r.Run()
}

func initDB() {
	var err error

	// 환경 변수에서 데이터베이스 접속 정보 읽기
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// dsn 생성
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, user, password, dbname, sslmode)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&domain.User{})
}
