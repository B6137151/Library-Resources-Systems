package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

const jwtSecret = "your_jwt_secret_here"

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password and store the user in the database
	user.Password = hashPassword(user.Password)
	// TODO: Store user in the database

	token, err := generateJWT(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func LoginUser(c *gin.Context) {
	var reqUser User
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: Retrieve the user from the database using the email
	var storedUser User
	// You will need to implement a function to retrieve the user from the database.

	if reqUser.Email != storedUser.Email || hashPassword(reqUser.Password) != storedUser.Password {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateJWT(&storedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

type Resource struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func CreateResource(c *gin.Context) {
	// Implement logic to create a resource in the database
	// ...

	c.JSON(201, gin.H{"message": "Resource created"})
}

func GetResource(c *gin.Context) {
	// Implement logic to get a resource from the database
	// ...

	resource := Resource{} // Replace with the fetched resource
	c.JSON(200, resource)
}

func GetResources(c *gin.Context) {
	// Implement logic to get all resources from the database
	// ...

	resources := []Resource{} // Replace with the fetched resources
	c.JSON(200, resources)
}
func UpdateResource(c *gin.Context) {
	// Implement logic to update a resource in the database
	// ...

	c.JSON(200, gin.H{"message": "Resource updated"})
}
func DeleteResource(c *gin.Context) {
	// Implement logic to delete a resource from the database
	// ...

	c.JSON(200, gin.H{"message": "Resource deleted"})
}
