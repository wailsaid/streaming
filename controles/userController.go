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

// SignUp registers a new user
// POST /api/auth/register
// Body JSON: { "username": "", "email": "", "password": "" }
func SignUp(c *CustomContext) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, Map{"error": "invalid request body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not hash password"})
		return
	}

	user := &models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}

	if err := database.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, Map{"error": "email already in use or could not create user"})
		return
	}

	log.Printf("new user registered: %s", user.Email)
	c.JSON(http.StatusCreated, Map{"message": "registered successfully"})
}

// Login authenticates a user and returns a JWT token as JSON
// POST /api/auth/login
// Body JSON: { "email": "", "password": "" }
func Login(c *CustomContext) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, Map{"error": "invalid request body"})
		return
	}

	user, err := database.FindUserByEmail(body.Email)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, Map{"error": "incorrect email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, Map{"error": "incorrect email or password"})
		return
	}

	// Generate JWT — 72 hour expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"uid": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	stringToken, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not generate token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    stringToken,
		Path:     "/",
		HttpOnly: true, // Prevents JavaScript access (Security!)
		Secure:   true, // Only sends over HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   72 * 3600, // 72 hours
	})

	c.JSON(http.StatusOK, Map{
		"token": stringToken,
		"user": Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// ListUsers returns all users (admin use)
func ListUsers(c *CustomContext) {
	users, err := database.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Map{"error": "could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Me returns the currently authenticated user from context
func Me(c *CustomContext) {
	// User is attached to context by JwtFilter
	// For now return a simple response — can be expanded
	c.JSON(http.StatusOK, Map{"message": "authenticated"})
}

// legacy HTML page funcs — kept as stubs so nothing breaks if referenced
func LoginPage(c *CustomContext)  { c.JSON(http.StatusOK, Map{"message": "use Vue login page"}) }
func SignupPage(c *CustomContext) { c.JSON(http.StatusOK, Map{"message": "use Vue register page"}) }
