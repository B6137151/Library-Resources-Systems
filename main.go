package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Replace username, password, and dbname with your MySQL credentials
	db, err := sql.Open("mysql", "library_app_user:4x@mpL3$pA$$w0rd@/library_resources")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// Define your API routes here
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Add user-related endpoints
	r.POST("/register", func(c *gin.Context) {
		registerUser(db, c)
	})
	r.POST("/login", func(c *gin.Context) {
		loginUser(db, c)
	})

	// Resource-related endpoints
	r.POST("/resources", createResource)
	r.GET("/resources/:id", getResource)
	r.GET("/resources", getResources)
	r.PUT("/resources/:id", updateResource)
	r.DELETE("/resources/:id", deleteResource)

	r.Run() // Listen and serve on 0.0.0.0:8080
}
