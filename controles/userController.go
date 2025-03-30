package controles

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

func SignUp(c *CustomContext) {
	log.Println(c.Params)

	var user models.User
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "faild to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "faild to hash password",
		})
		return
	}
	user.Password = string(hash)

	err = database.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "could not store user",
		})
		return
	}

	c.Redirect(302, "/login")
}

func ListUsers(c *CustomContext) {
	ListUsers, err := database.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "could not retrieve users",
		})
		return
	}
	c.JSON(http.StatusOK, ListUsers)
}

func Login(c *CustomContext) {
	var reqBody struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	if c.Bind(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "could not read request body",
		})
		return
	}

	user, err := database.FindUserByEmail(reqBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "incorrect email or password",
		})
		return
	}

	// c.JSON()

	res := database.DB.First(&user, "email = ?", reqBody.Email)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, Map{
			"error": "incorrect email or password ",
		})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, Map{
			"error": "incorrect email or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, Map{
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
		c.JSON(http.StatusBadRequest, Map{
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

func LoginPage(c *CustomContext) {
	log.Println("login.html")
	c.HTML(200, "login", Map{
		"title": "Login",
	})
}

func SignupPage(c *CustomContext) {
	log.Println("signup.html")
	c.HTML(200, "signup", Map{
		"title": "Sign Up",
	})
}
