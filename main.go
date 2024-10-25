package main

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
	"be-bwastartup/docs"
	"be-bwastartup/handler"
	"be-bwastartup/helper"
	"be-bwastartup/payment"
	"be-bwastartup/transaction"
	"be-bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "be-bwastartup/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Crowdfunding API
// @version 1.0
// @description Crowdfunding API Description

func main(){

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
	paymentService := payment.NewService()
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Use(cors.Default())
	api := router.Group("/api/v1")
	api.Use(AuthMiddleware(authService, userService))

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
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	
	// transaction
	api.GET("/transactions/campaign/:campaign_id", transactionHandler.GetTransactionsByCampaignID)
	api.GET("/transactions", transactionHandler.GetTransactionByUserID)
	api.POST("/transactions", transactionHandler.CreateTransaction)
	router.POST("/api/v1/transactions/notification", transactionHandler.GetNotification)


	router.GET("/", func(ctx *gin.Context) {
		response := helper.APIResponse("APP IS ONLINE", http.StatusOK, "Success", nil)
		ctx.JSON(http.StatusOK, response)
	})
	

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost" + os.Getenv("PORT"))
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	router.Run(os.Getenv("PORT"))

}

func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context){
		autHeader	:= c.GetHeader("Authorization")
		if !strings.Contains(autHeader, "Bearer"){
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		t := strings.Split(autHeader, " ")
		if len(t) != 2 {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		validatedToken, err := authService.ValidateToken(t[1])
		if err != nil {
			errorsResponse := gin.H{"errors": err.Error()}
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", errorsResponse)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := validatedToken.Claims.(jwt.MapClaims)
		if !ok || !validatedToken.Valid {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)	
	}
}