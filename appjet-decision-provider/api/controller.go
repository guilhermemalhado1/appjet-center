package api

import (
	"appjet-decision-provider/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	services.LoginHandler(c)
}

func SignupHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if services.CheckIfTokenValid(token) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Invalid"})
		return
	}

	services.CreateUserHandler(c)
}

func LogoutHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if services.CheckIfTokenValid(token) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Invalid"})
		return
	}

	services.LogoutHandler(c)
}

// GenericHandler is a generic handler function for endpoints other than login, logout, and signup
func GenericHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if services.CheckIfTokenValid(token) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Invalid"})
		return
	}

	services.GenericHandler(c)

}
