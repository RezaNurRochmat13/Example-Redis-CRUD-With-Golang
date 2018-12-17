package controller

import (
	"encoding/json"
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

	// Parsing data with object property
	parseRecord, err := json.Marshal(createBooksPayload)

	// Saving into Redis
	_, err = connection.Do("SET", UUIDBooks, parseRecord)

	// Exception error handling
	exception.GlobalException(err)

	// Message when trigerred successfully
	c.JSON(http.StatusOK, gin.H{"message": "Inserted successfully"})

	// Closing connection Redis when completed
	defer connection.Close()

}

// GetDetailBooksFromRedis func get detail books
func GetDetailBooksFromRedis(c *gin.Context) {
	// Initialize redis connection
	connection := pool.Get()

	// Declare models
	var Books model.Books

	// Set parameters
	UUIDBooks := c.Param("UUIDBooks")

	// Get values from Redis
	value, err := redis.String(connection.Do("GET", UUIDBooks))

	// Unmarshalling values and map to struct from Redis
	err = json.Unmarshal([]byte(value), &Books)

	if err != nil {
		panic(err)
	}

	// Serve as JSON formats
	c.JSON(http.StatusOK, gin.H{
		"data": Books,
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
		exception.GlobalException(err)

		// Closing connection Redis when completed
		defer connection.Close()
	}

}

// DeleteBooksFromRedisAndDB func does delete data from redis and DB
func DeleteBooksFromRedisAndDB(context *gin.Context) {

	// Initialize redis connection
	connection := pool.Get()

	// Initialize db connection
	db := config.DatabaseConn()

	// Declare model
	var deleteBooks model.Books

	// Set parameter UUID
	UUIDBooks := context.Param("UUIDBooks")

	// Check record in database
	if db.Where("books.uuid_books = ?", UUIDBooks).Find(&deleteBooks).RecordNotFound() {

		// When not found record
		context.JSON(http.StatusOK, gin.H{"message": "Record not found"})

		// Closing database when completed execution
		defer db.Close()

	} else {

		// When found record, deleted record in database
		db.Where("books.uuid_books = ?", UUIDBooks).Find(&deleteBooks).Delete(&deleteBooks)

		// When record found, deleted too in Redis
		_, err := connection.Do("DEL", UUIDBooks)

		// Exception handling
		exception.GlobalException(err)

		// When successfully deleted
		context.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})

		// Closing database when completed execution
		defer db.Close()

		// Closing redis connection when compeleted execution
		defer connection.Close()
	}
}
