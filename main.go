package main

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
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

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)
	log.SetOutput(os.Stdout)
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connection Success")



	PORT := os.Getenv("PORT")

	
	// repositorys
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// services
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, paymentService, campaignRepository)

	// handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService, authService)
	transactionHandler := handler.NewTransactionHandler(transactionService, campaignService, paymentService)

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gin.Logger())
	
	api := router.Group("/api/v1")
	api.Use(authMiddleware(authService, userService))

	router.Static("/images", "./images")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.GET("/users", userHandler.GetUsers)
	
	// users
	api.POST("/avatar", userHandler.UploadAvatar)
	api.GET("/me", userHandler.GetUser)

	// campaign
	api.GET("/campaigns", campaignHandler.FindCampaigns)
	api.GET("/campaigns/:id", campaignHandler.FindCampaign)
	api.POST("/campaigns", campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	api.POST("/campaigns/images", campaignHandler.UploadImage)

	// transaction
	api.GET("/campaigns/:id/transactions", transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", transactionHandler.GetUserTransactions)
	api.POST("/transactions", transactionHandler.CreateTransaction)
	router.POST("/api/v1/transactions/notifications", transactionHandler.GetNotification)

	router.Run(PORT)

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer"){
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			
		}

		tokenPayload := strings.Split(authHeader, " ")[1]

		token, err := authService.ValidateToken(tokenPayload)
		if err != nil {
			errors := gin.H{"errors": err.Error()}
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			errors := gin.H{"errors": err.Error()}
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}


