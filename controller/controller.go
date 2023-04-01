package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SetDBtoContext(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

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
	UserID      int    `json:"user_id"`
}

// Implement the storeUser function
func storeUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (email, password, role) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.Email, user.Password, user.Role)
	return err
}
func getUserByEmail(db *sql.DB, email string) (*User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = ?"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func storeResource(db *sql.DB, resource *Resource) (int64, error) {
	query := "INSERT INTO resources (title, author, description, user_id) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, resource.Title, resource.Author, resource.Description, resource.UserID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// Implement the getResourceByID function
func getResourceByID(db *sql.DB, id int) (*Resource, error) {
	query := "SELECT id, title, author, description, user_id FROM resources WHERE id = ?"
	row := db.QueryRow(query, id)

	var resource Resource
	err := row.Scan(&resource.ID, &resource.Title, &resource.Author, &resource.Description, &resource.UserID)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Implement the getAllResources function
func getAllResources(db *sql.DB) ([]Resource, error) {
	query := "SELECT id, title, author, description, user_id FROM resources"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var resource Resource
		err := rows.Scan(&resource.ID, &resource.Title, &resource.Author, &resource.Description, &resource.UserID)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func CreateResource(c *gin.Context) {
	var resource Resource
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: Store the resource in the database using the storeResource function
	// ...

	c.JSON(201, gin.H{"message": "Resource created"})
}

func GetResource(c *gin.Context) {
	// Get the database connection from the context
	db, exists := c.Get("db")
	if !exists {
		c.JSON(500, gin.H{"error": "Database connection not found"})
		return
	}

	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	// Retrieve the resource from the database using the getResourceByID function
	resource, err := getResourceByID(db.(*sql.DB), resourceID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error retrieving resource"})
		return
	}

	c.JSON(200, resource)
}

func GetResources(c *gin.Context) {
	// TODO: Retrieve all resources from the database using the getAllResources function
	// ...

	resources := []Resource{} // Replace with the fetched resources
	c.JSON(200, resources)
}
func updateResource(db *sql.DB, resource *Resource) error {
	query := "UPDATE resources SET title = ?, author = ?, description = ?, user_id = ? WHERE id = ?"
	_, err := db.Exec(query, resource.Title, resource.Author, resource.Description, resource.UserID, resource.ID)
	return err
}
func UpdateResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var resource Resource
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	resource.ID = id

	// Get the database connection from the context
	db, exists := c.Get("db")
	if !exists {
		c.JSON(500, gin.H{"error": "Database connection not found"})
		return
	}

	// Update the resource in the database using the updateResource function
	err = updateResource(db.(*sql.DB), &resource)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating resource"})
		return
	}

	c.JSON(200, gin.H{"message": "Resource updated"})
}

func DeleteResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	// Get the database connection from the context
	db, exists := c.Get("db")
	if !exists {
		c.JSON(500, gin.H{"error": "Database connection not found"})
		return
	}

	err = deleteResource(db.(*sql.DB), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting resource"})
		return
	}

	c.JSON(200, gin.H{"message": "Resource deleted"})
}

func deleteResource(db *sql.DB, id int) error {
	query := "DELETE FROM resources WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}
