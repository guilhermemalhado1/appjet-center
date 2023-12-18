// main.go
package main

import (
	"appjet-decision-provider/api"
	"appjet-decision-provider/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	services.LoadConfig()
	services.InitDB()

	// Define API routes
	apiGroup := r.Group("/api")
	{
		//do login
		apiGroup.POST("/login", api.LoginHandler)
		//do logout
		apiGroup.GET("/logout", api.LogoutHandler)
		//create user
		apiGroup.POST("/signup", api.SignupHandler)
		// Generic handler for other endpoints
		apiGroup.Use(api.GenericHandler)
	}

	// Run the server
	r.Run(":8080")
}
