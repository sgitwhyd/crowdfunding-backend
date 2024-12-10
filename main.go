package main

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
	clod "be-bwastartup/cloudinary"
	"be-bwastartup/handler"
	"be-bwastartup/helper"
	"be-bwastartup/payment"
	"be-bwastartup/redis"
	"be-bwastartup/transaction"
	"be-bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"be-bwastartup/docs"

	"be-bwastartup/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Crowdfunding API
// @version 1.0
// @description Crowdfunding API Description
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization



func main(){

	time.LoadLocation("Asia/Jakarta")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print(err)
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("GIN_MODE") == "release" {	
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	dsn := os.Getenv("POSTGRES_URL")
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection Success")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// services

	cloudinaryService, err := clod.NewService()
	if err != nil {
		log.Fatal(err.Error())
	}
	paymentService := payment.NewService()
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository, cloudinaryService)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	redisService := redis.NewService()
	userService := user.NewService(userRepository, cloudinaryService)
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService, redisService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Use(cors.Default())
	api := router.Group("/api/v1")
	router.Static("/images", "./images")
	api.Use(middleware.Auth(authService, userService, redisService))

	router.POST("/api/v1/sessions", userHandler.Login)

	router.POST("/api/v1/users", userHandler.RegisterUser)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)
	api.PUT("/users", userHandler.UpdateUser)
	api.GET("/users/current", userHandler.GetCurrentUser)

	// campaign
	api.POST("/campaigns", campaignHandler.CreateCampaign)
	router.GET("/api/v1/campaigns", campaignHandler.GetCampaigns)
	api.POST("/campaigns/images", campaignHandler.SaveCampaignImage)
	router.GET("/api/v1/campaigns/:id", campaignHandler.GetCampaign)
	api.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	
	// transaction
	api.GET("/transactions/campaign/:campaign_id", transactionHandler.GetTransactionsByCampaignID)
	api.GET("/transactions", transactionHandler.GetTransactionByUserID)
	api.POST("/transactions", transactionHandler.CreateTransaction)
	router.POST("/api/v1/transactions/notification", transactionHandler.GetNotification)


	router.GET("/", func(ctx *gin.Context) {
		response := helper.APIResponse("API IS ONLINE", http.StatusOK, "Success", nil)
		ctx.JSON(http.StatusOK, response)
	})
	

	docs.SwaggerInfo.Host = fmt.Sprint(os.Getenv("BASE_URL"))
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	router.Run(os.Getenv("PORT"))

}