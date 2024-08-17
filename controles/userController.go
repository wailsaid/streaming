package controles

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/saidwail/streaming/initEnv"
	"github.com/saidwail/streaming/models"
)

func SignUp(c *gin.Context) {
	log.Println(c.Params)

	var user models.User
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "faild to read body",
		})
		return
	}
	/* log.Println(user.Email)
	log.Println(user.User)
	log.Println(user.Password) */

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "faild to hash password",
		})
		return
	}
	user.Password = string(hash)

	res := initEnv.DB.Create(&user)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not store user",
		})
		return
	}

	c.Redirect(302, "/login")
}

func ListUsers(c *gin.Context) {
	var ListUsers []models.User
	res := initEnv.DB.Find(&ListUsers)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retreve users",
		})
		return
	}
	c.JSON(http.StatusOK, ListUsers)
}

func Login(c *gin.Context) {
	var reqBody struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	if c.Bind(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not read request body",
		})
		return
	}

	// c.JSON()
	var user models.User
	res := initEnv.DB.First(&user, "email = ?", reqBody.Email)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect email or password ",
		})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect email or password",
		})
		return
	}
	// c.JSON(200, gin.H{})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	})

	stringToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "faild to generate token " + err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		stringToken,
		int(time.Now().Add(time.Minute*1).Unix()),
		"",
		"",
		true,
		true,
	)
	c.Redirect(302, "/")
}
