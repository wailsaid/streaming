package midelware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/saidwail/streaming/initEnv"
	"github.com/saidwail/streaming/models"
)

func JwtFilter(c *gin.Context) {
	stringToken, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["sub"])
		var user models.User

		res := initEnv.DB.First(&user, "email = ?", claims["sub"])

		if res.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{})
		}

		c.Next()
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
