package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
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

func storeUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (email, password, role) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.Email, user.Password, user.Role)
	return err
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = ?"
	row := db.QueryRow(query, email)

	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func registerUser(db *sql.DB, c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user.Password = hashPassword(user.Password)

	if err := storeUser(db, &user); err != nil {
		c.JSON(500, gin.H{"error": "Error storing user"})
		return
	}

	token, err := generateJWT(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func loginUser(db *sql.DB, c *gin.Context) {
	var reqUser User
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	storedUser, err := getUserByEmail(db, reqUser.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	if hashPassword(reqUser.Password) != storedUser.Password {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateJWT(storedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating JWT"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
