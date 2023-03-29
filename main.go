package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Replace username, password, and dbname with your MySQL credentials
	db, err := sql.Open("mysql", "username:password@/dbname")
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

	r.Run() // Listen and serve on 0.0.0.0:8080
}
