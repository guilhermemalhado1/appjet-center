package api

import (
	"appjet-identity-provider/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func verifyIfIsAlreadyLoggedIn(token string) bool {
	return true;
}

func LoginHandler(c *gin.Context) {

	services.LoginHandler(c)
}

func SignupHandler(c *gin.Context) {

	services.CreateUserHandler(c)

}

func LoadLoginStateHandler(c *gin.Context) {
	// Implement your logic to check the login state
	// ...

	// Return success or failure response
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Login state retrieved"})
}

func DeleteLoginStateHandler(c *gin.Context) {
	// Implement your logic to delete the login state
	// ...

	// Return success or failure response
	c.JSON(http.StatusOK, gin.H{"message": "Login state deleted"})
}
