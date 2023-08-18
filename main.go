package main

import (
	"go-gin-postgre/domain"
	"go-gin-postgre/domain/usecases"
	"go-gin-postgre/infrastructures"
	"go-gin-postgre/interfaces/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	initDB()

	// 의존성 주입
	repo := infrastructures.NewUserRepository(db)
	usecase := usecases.NewUserUsecase(repo)
	handler := handlers.NewUserHandler(usecase)

	r := gin.Default()
	r.GET("/users", handler.GetUsers)
	r.GET("/users/:id", handler.GetUser)
	r.POST("/users", handler.CreateUser)
	r.PUT("/users/:id", handler.UpdateUser)
	r.DELETE("/users/:id", handler.DeleteUser)

	r.Run()
}

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&domain.User{})
}
