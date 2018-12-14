package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"go-example-redis/config"
	"go-example-redis/exception"
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

	payloads, err := json.Marshal(createBooksPayload)

	exception.CauseException(err)

	log.Println("JSON data", payloads)

	// Saves into DB
	db.Create(&createBooksPayload)

	// Saving into Redis
	_, err = connection.Do("SET", UUIDBooks, createBooksPayload)

	// Exception error handling
	exception.CauseException(err)

	// Message when trigerred successfully
	c.JSON(http.StatusOK, gin.H{"message": "Inserted successfully"})

	// Closing connection Redis when completed
	defer connection.Close()

}

func GetDetailBooksFromRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Set parameters
	UUIDBooks := c.Param("UUIDBooks")

	// Get values from Redis
	value, err := redis.String(connection.Do("GET", UUIDBooks))

	// Exception error handling
	exception.CauseException(err)

	//// Serve as JSON formats
	c.JSON(http.StatusOK, gin.H{
		"data": value,
		"keys": UUIDBooks})

	// Closing connection Redis when successfully complete
	defer connection.Close()
}

func UpdateBooksInRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	db := config.DatabaseConn()

	var updateBooks model.Books

	// Set parameters
	UUIDBooks := c.Param("UUIDBooks")
	BooksName := c.Param("BooksName")
	BooksPublisher := c.Param("BooksPublisher")
	BooksWriter := c.Param("BooksWriter")
	BooksDescription := c.Param("BooksDescription")

	updateBookPayload := model.Books{BooksName: BooksName, BooksPublisher: BooksPublisher,
		BooksWriter: BooksWriter, BooksDescription: BooksDescription}

	c.BindJSON(&updateBookPayload)

	if db.Where("books.uuid_books = ?", UUIDBooks).Find(&updateBooks).RecordNotFound() {
		c.JSON(http.StatusOK, gin.H{"message": "Record not found"})
	} else {
		db.Where("products.uuid_products = ?", UUIDBooks).Table("products").Update(&updateBooks)

		_, err := connection.Do("SET", UUIDBooks, updateBookPayload)

		// Exception error handling
		exception.CauseException(err)

		c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})

	}

}
