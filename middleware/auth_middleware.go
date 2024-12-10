package middleware

import (
	"be-bwastartup/auth"
	"be-bwastartup/helper"
	rdb "be-bwastartup/redis"
	"be-bwastartup/user"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)


func Auth(authService auth.Service, userService user.Service, rdb rdb.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		autHeader := c.GetHeader("Authorization")
		if !strings.Contains(autHeader, "Bearer") {
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

		sessionKey := "session" + t[1]
		sessionData, err := rdb.Get(sessionKey)
		if err != redis.Nil {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var user user.User
		err = json.Unmarshal([]byte(sessionData), &user)
		if err != nil {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		} else {
			c.Set("currentUser", user)
			c.Next()
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
		u, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userJSON, err := json.Marshal(u)
		if err != nil {
			response := helper.APIResponse("Unauthorize", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		_,err = rdb.Save(sessionKey, userJSON, time.Hour)
		c.Set("currentUser", u)
		c.Next()
	}
}