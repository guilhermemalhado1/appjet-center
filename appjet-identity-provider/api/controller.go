package api

import (
	"appjet-identity-provider/services"
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
