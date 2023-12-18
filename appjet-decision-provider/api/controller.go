package api

import (
	"appjet-decision-provider/services"
	"bytes"
	"io/ioutil"
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
	// Get the incoming request details
	token := c.GetHeader("Authorization")
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// Check if the token is valid
	if !services.CheckIfTokenValid(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Invalid"})
		return
	}

	// Make an HTTP request to another URL with the same attributes
	targetURL := "https://example.com/target" // Replace with your target URL
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers from the incoming request to the new request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make HTTP request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Set the same status code and response body as the target URL's response
	c.JSON(resp.StatusCode, gin.H{"message": "Generic handler for other endpoints", "target_response": string(responseBody)})
}
