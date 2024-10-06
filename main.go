package main

import (
	"be-bwastartup/auth"
	"be-bwastartup/handler"
	"be-bwastartup/user"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main(){
	gin.SetMode(gin.DebugMode)
	log.SetOutput(os.Stdout)
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connection Success")
	// repositorys
	userRepository := user.NewRepository(db)

	// services
	userService := user.NewService(userRepository)
	authService := auth.NewService()


	userHandler := handler.NewUserHandler(userService, authService)


	router := gin.Default()
	router.Use(gin.Logger())
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.GET("/users", userHandler.GetUsers)
	api.POST("/avatar", userHandler.UploadAvatar)


	router.Run()

}
