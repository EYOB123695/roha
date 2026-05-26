package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/EYOB123695/roha/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("In middleware")

		// Get the cookie of the request 
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check expiration
			if exp, ok := claims["exp"].(float64); ok {
				if float64(time.Now().Unix()) > exp {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Find user
			userID := uint(claims["sub"].(float64))
			user, err := userRepo.GetByID(userID)
			if err != nil || user == nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Attach user to context
			c.Set("user", *user)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}