package midelware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/saidwail/streaming/initEnv"
	"github.com/saidwail/streaming/models"
)

func JwtFilter(c *gin.Context) {
	/* if c.Request.URL.Path == "/" {
		c.Next()
		return

	} */
	unauthorized := false

	stringToken, err := c.Cookie("Authorization")
	if err != nil {
		unauthorized = true
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
		if c.Request.URL.Path == "/" {
			c.Next()
			return
		}
	}

	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		unauthorized = true
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["sub"])
		var user models.User

		res := initEnv.DB.First(&user, "email = ?", claims["sub"])

		if res.Error != nil {
			unauthorized = true
			//c.JSON(http.StatusUnauthorized, gin.H{})
			//return
		}
		if user.ID == 0 {
			unauthorized = true
			//c.JSON(http.StatusUnauthorized, gin.H{})
			//return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			unauthorized = true
			//c.JSON(http.StatusUnauthorized, gin.H{})
		}
		c.Set("logged_in", true)
		c.Next()
		return
	}
	if unauthorized {
		c.Redirect(302, "/login")
	}
	//c.AbortWithStatus(http.StatusUnauthorized)
}
