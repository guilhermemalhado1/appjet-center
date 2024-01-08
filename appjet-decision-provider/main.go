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

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.GET("/logout", api.LogoutHandler)
		apiGroup.POST("/signup", api.SignupHandler)
		apiGroup.Any("/:path", api.GenericHandler)
	}

	r.Run(":8080")
}
