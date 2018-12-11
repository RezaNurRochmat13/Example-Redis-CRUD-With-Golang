package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"go-example-redis/config"
	"go-example-redis/model"
	"log"
	"net/http"
)

// Declare global configuration database pooling Redis
var pool = config.NewConnectionPools()

// SetBooksValuesToDBAndRedis func does set data to DB and Redis
func CreateBooksToDBAndRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Closing connection Redis when idle
	defer connection.Close()

	// Initialize database connection
	db := config.DatabaseConn()

	// Set parameters books
	UUIDBooks := uuid.Must(uuid.NewV4())
	BooksName := c.Param("BooksName")
	BooksPublisher := c.Param("BooksPublisher")
	BooksWriter := c.Param("BooksWriter")
	BooksDescription := c.Param("BooksDescription")

	// Create books payload insert
	createBooksPayload := model.Books{UUIDBooks: UUIDBooks, BooksName: BooksName, BooksPublisher: BooksPublisher,
		BooksWriter: BooksWriter, BooksDescription: BooksDescription}

	// Bind as JSON format
	c.BindJSON(&createBooksPayload)

	// Saves into DB
	db.Create(&createBooksPayload)

	// Saving into Redis
	_, err := connection.Do("SET", UUIDBooks, createBooksPayload)

	// Exception error handling
	if err != nil {
		panic(err.Error())
	}

	// Message when trigerred successfully
	c.JSON(http.StatusOK, gin.H{"message": "Inserted successfully"})

}

func GetDetailBooksFromRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Closing connection Redis when idle
	defer connection.Close()

	// Set parameters
	UUIDBooks := c.Param("UUIDBooks")

	// Get values from Redis
	value, err := redis.String(connection.Do("GET", UUIDBooks))

	log.Println("Books value from Redis", value)

	// Exception error handling
	if err != nil {
		panic(err)
	}

	//// Serve as JSON formats
	c.JSON(http.StatusOK, gin.H{"data": value})
}
