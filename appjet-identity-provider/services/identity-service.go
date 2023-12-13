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

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
