package main

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
	"be-bwastartup/handler"
	"be-bwastartup/helper"
	"be-bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main(){
	gin.SetMode(gin.DebugMode)
	dsn := "postgres://default:UErH05WiceYt@ep-white-rain-a1j4lx8q.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require"
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection Success")
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	// services
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)


	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(AuthMiddleware(authService, userService))


	router.POST("/api/v1/users", userHandler.RegisterUser)
	router.POST("/api/v1/sessions", userHandler.Login)
	api.POST("/email_checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)
	api.PUT("/users", userHandler.UpdateUser)

	// campaign

	api.POST("/campaigns", campaignHandler.CreateCampaign)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.POST("/campaigns/images", campaignHandler.SaveCampaignImage)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)

	router.Run()

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