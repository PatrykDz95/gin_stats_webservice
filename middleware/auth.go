package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"gin/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
	"net/http"
	"time"
)

const pepper = "pepper"

type User struct {
	Username string
	Password string
	Email    string
	Salt     string
	Verified bool
	Token    string
}

func ChangePasswordHandler(c *gin.Context) {
	var credentials struct {
		Username    string
		OldPassword string
		NewPassword string
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve user from database
	var user User
	if userLogin := database.DB.Where("username = ?", credentials.Username).First(&user); userLogin.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare input password with hashed password from the database
	hashedInputPassword := hashPassword(credentials.OldPassword, user.Salt, pepper)

	// Compare the hashed input password with the stored hashed password
	if user.Password != hashedInputPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Hash the new password
	salt, err := generateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}
	hashedPassword := hashPassword(credentials.NewPassword, salt, pepper)

	// Update the user's password and salt in the database
	user.Password = hashedPassword
	user.Salt = salt
	userUpdated := database.DB.Where("username = ?", credentials.Username).Updates(&user)
	if userUpdated.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": userUpdated.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})

}

func RegisterHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password
	salt, err := generateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}
	hashedPassword := hashPassword(user.Password, salt, pepper)
	user.Password = hashedPassword
	user.Salt = salt

	// Store the user in the database
	userCreated := database.DB.Create(&user)
	if userCreated.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": userCreated.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func hashPassword(password, salt string, pepper string) string {
	// Convert the salt and pepper to byte slices
	saltBytes, _ := base64.RawStdEncoding.DecodeString(salt)
	pepperBytes := []byte(pepper)

	// Combine the password and pepper
	passwordAndPepper := password + string(pepperBytes)

	// Hash the password, salt and pepper using Argon2
	hash := argon2.IDKey([]byte(passwordAndPepper), saltBytes, 1, 64*1024, 4, 32)

	return base64.RawStdEncoding.EncodeToString(hash)
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func LoginHandler(c *gin.Context) {
	var credentials struct {
		Username string
		Password string
		Email    string
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve user from database
	var user User
	if userLogin := database.DB.Where("username = ?", credentials.Username).First(&user); userLogin.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare input password with hashed password from the database

	hashedInputPassword := hashPassword(credentials.Password, user.Salt, pepper)

	// Compare the hashed input password with the stored hashed password
	if user.Password != hashedInputPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and return access token
	accessToken, err := generateAccessToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract access token from request header, query parameter, or cookie
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate access token
		validated, err := validateAccessToken(accessToken)
		if !validated || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token," + err.Error()})
			c.Abort()
			return
		}

		// If token is valid, continue to the next middleware/handler
		c.Next()
	}
}

func generateAccessToken(username string) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the username and expiration time
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// Generate JWT token with the claims and sign it with a secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateAccessToken(accessToken string) (bool, error) {
	// Parse the JWT token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Provide the secret key to verify the token
	})
	if err != nil {
		return false, err
	}

	// Check if the token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}
