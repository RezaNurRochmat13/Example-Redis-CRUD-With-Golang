package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"go-example-redis/config"
	"go-example-redis/exception"
	"go-example-redis/model"
	"net/http"
)

// Declare global configuration database pooling Redis
var pool = config.NewConnectionPools()

// SetBooksValuesToDBAndRedis func does set data to DB and Redis
func CreateBooksToDBAndRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

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
	exception.GlobalMessageException(err)

	// Message when trigerred successfully
	c.JSON(http.StatusOK, gin.H{"message": "Inserted successfully"})

	// Closing connection Redis when completed
	defer connection.Close()

}

// GetDetailBooksFromRedis func get detail books
func GetDetailBooksFromRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Set parameters
	UUIDBooks := c.Param("UUIDBooks")

	// Get values from Redis
	value, err := redis.String(connection.Do("GET", UUIDBooks))

	// Exception error handling
	exception.GlobalMessageException(err)

	//// Serve as JSON formats
	c.JSON(http.StatusOK, gin.H{
		"data": value,
		"keys": UUIDBooks})

	// Closing connection Redis when successfully complete
	defer connection.Close()
}

// Update books data in database and redis
func UpdateBooksInRedisAndDB(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Initialize database connection
	db := config.DatabaseConn()

	// Declare model books
	var updateBooks model.Books

	// Set parameters input
	UUIDBooks := c.Param("UUIDBooks")
	BooksName := c.Param("BooksName")
	BooksPublisher := c.Param("BooksPublisher")
	BooksWriter := c.Param("BooksWriter")
	BooksDescription := c.Param("BooksDescription")

	// Create a payload update
	updateBookPayload := model.Books{BooksName: BooksName, BooksPublisher: BooksPublisher,
		BooksWriter: BooksWriter, BooksDescription: BooksDescription}

	// Bind as JSON formats
	c.BindJSON(&updateBookPayload)

	// Check record in database is exist or not
	if db.Where("books.uuid_books = ?", UUIDBooks).Find(&updateBooks).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{"message": "Record not found"})

		// Closing database connection
		defer db.Close()
	} else {
		db.Where("books.uuid_books = ?", UUIDBooks).Table("books").Update(&updateBookPayload)

		c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})

		// Closing database connection
		defer db.Close()

		// Update record in Redis Database
		_, err := connection.Do("SET", UUIDBooks, updateBookPayload)

		// Exception error handling
		exception.GlobalMessageException(err)

		// Closing connection Redis when completed
		defer connection.Close()
	}

}
