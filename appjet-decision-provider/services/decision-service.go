package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const jsonPropertyFileName = "config.json"

// User represents a user in the system.
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type AppjetConfiguration struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user_id"` // You may need to adjust the data type based on your requirements
	Config   string `db:"config"`
}

type AppjetConfigurationDto struct {
	AppJetURL string `json:"appJetURL"`
	Plugins   struct {
		Git struct {
			Enabled      bool   `json:"enabled"`
			RepoURL      string `json:"repo-url"`
			RepoUser     string `json:"repo-user"`
			RepoPassword string `json:"repo-password"`
		} `json:"git"`
	} `json:"plugins"`
	Cluster struct {
		Name    string `json:"name"`
		Servers []struct {
			Name      string `json:"name"`
			IPAddress string `json:"ip-address"`
			Port      int    `json:"port"`
			User      string `json:"user"`
			Password  string `json:"password"`
		} `json:"servers"`
	} `json:"cluster"`
	Artifact struct {
		Language struct {
			Name      string `json:"name"`
			Version   string `json:"version"`
			Framework struct {
				Name string `json:"name"`
			} `json:"framework"`
		} `json:"language"`
	} `json:"artifact"`
}

type Token struct {
	UserID  uint   // Foreign key referencing users(id)
	Token   string `json:"token"`
	Expired bool   `json:"expired"`
}

// Config represents the configuration structure.
type Config struct {
	DBDriver           string `json:"DBDriver"`
	DBConnectionString string `json:"DBConnectionString"`
}

var db *sqlx.DB

var AppConfig Config

func LoadConfig() {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		panic(err)
	}
}

func getConfigFilePath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get current file path")
	}

	return filepath.Join(filepath.Dir(filename), "", jsonPropertyFileName), nil
}

func InitDB() {
	var err error

	db, err = sqlx.Open(AppConfig.DBDriver, AppConfig.DBConnectionString)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Ping the database to verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	fmt.Println("Connected to the database")
}

func CreateUserHandler(c *gin.Context) {

	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Insert user into the database
	_, err := db.NamedExec(`
		INSERT INTO users (username, email, password) VALUES (:username, :email, :password)
	`, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func LoginHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Query the database for the user
	err := db.Get(&user, `
		SELECT * FROM users WHERE username = $1 AND password = $2
	`, user.Username, user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		}
		return
	}

	token := uuid.New()

	// Insert the new token into the user_token table
	// Insert the new token into the user_token table using NamedExec
	_, err = db.NamedExec(`
		INSERT INTO user_token (user_id, token, expired) VALUES (:user_id, :token, :expired)
	`, map[string]interface{}{
		"user_id": user.ID,
		"token":   token.String(),
		"expired": false,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ACCESS_TOKEN": token.String()})
}

func CheckIfTokenValid(token string) bool {
	var count int

	// Query the database for the count of rows
	err := db.Get(&count, `
		SELECT COUNT(*) FROM user_token WHERE token = $1
	`, token)

	if err != nil {
		return false
	}

	// Check if exactly 1 row is returned
	if count != 1 {
		return false
	}

	return true
}

func GenericHandler(c *gin.Context) {
	// Get the requested endpoint (path)
	endpointRequested := strings.TrimPrefix(c.Request.URL.Path, "/api/")
	fmt.Println("Endpoint Requested:", endpointRequested)

	switch endpointRequested {
	case "start":
		StartHandler(c)
	case "ls":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	case "inspect":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	case "show":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	case "rename":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	case "delete":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	case "create":
		c.JSON(http.StatusAccepted, gin.H{"message": "IN PROGRESS"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Operation not suported"})
		return
	}
}

func LogoutHandler(c *gin.Context) {
	// Retrieve the token from the request header
	token := c.GetHeader("Authorization")

	// Check if the token is empty
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	// Delete the row from the user_token table where the token matches
	result, err := db.Exec(`
		DELETE FROM user_token WHERE token = $1
	`, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Check if the row was deleted
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	return
}

func StartHandler(c *gin.Context) {
	var startRequest AppjetConfigurationDto

	if err := c.ShouldBindJSON(&startRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Printf("Received Start Request:\n%+v\n", startRequest)

	// Retrieve the token from the request header
	token := c.GetHeader("Authorization")

	// Check if the token is empty
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	// Query the database to get the user ID based on the token
	var userID int
	err := db.Get(&userID, `
		SELECT user_id FROM user_token WHERE token = $1
	`, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID from token"})
		return
	}

	// Convert startRequest to JSON string
	configJSON, err := json.Marshal(startRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process start request"})
		return
	}

	// Insert data into the database
	_, err = db.NamedExec(`
		INSERT INTO appjet_configurations (user_id, config) VALUES (:user_id, :config)
	`, map[string]interface{}{
		"user_id": userID,
		"config":  string(configJSON),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert start request into the database"})
		return
	}

	// Respond with the full JSON
	c.JSON(http.StatusOK, gin.H{"message": "Start request processed successfully", "config": startRequest})
}


