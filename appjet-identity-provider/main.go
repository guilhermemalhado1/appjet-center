// main.go
package main

import (
	"appjet-identity-provider/api"
	"appjet-identity-provider/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	services.LoadConfig()
	services.InitDB()

	// Define API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.GET("/login", api.LoadLoginStateHandler)
		apiGroup.GET("/logout", api.DeleteLoginStateHandler)
		apiGroup.POST("/signup", api.SignupHandler)
	}

	// Run the server
	r.Run(":8080")
}
