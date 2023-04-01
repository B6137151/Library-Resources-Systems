package main

import (
	"database/sql"
	"github.com/B6137151/library-resources-system/controller"
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

	r.Use(controller.SetDBtoContext(db))

	// Define your API routes here
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "โจ๋",
		})
	})

	// Add user-related endpoints
	r.POST("/register", controller.RegisterUser)
	r.POST("/login", controller.LoginUser)

	// Resource-related endpoints
	r.POST("/resources", controller.CreateResource)
	r.GET("/resources/:id", controller.GetResource)
	r.GET("/resources", controller.GetResources)
	r.PUT("/resources/:id", controller.UpdateResource)
	r.DELETE("/resources/:id", controller.DeleteResource)

	r.Run() // Listen and serve on 0.0.0.0:8080
}
