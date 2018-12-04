package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-redis/model"
)

// SetDataBooksIntoRedis func does set value into Redis
func SetDataBooksIntoRedis(cmd redis.Conn) error {
	// Set struct values
	Products := &model.Books{
		UUIDBooks:       "1e9ff5c4-7f96-47f7-b1ab-bb995bc15041",
		BookName:        "Harry potter",
		BookPublisher:   "Gramedia Pub",
		BookWriters:     "J.K Rowling",
		BookDescription: "Harry potter books",
	}

	// Encode struct to JSON using Marshal
	jsonEncode, err := json.Marshal(Products)

	// Exception error handling
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Serialized data books : ", string(jsonEncode))
	}

	// Set values into Redis
	_, err = cmd.Do("SET", "Book-1", jsonEncode)

	// Exception error handling
	if err != nil {
		panic(err.Error())
	}

	// Returned nil / return 0
	return nil
}

// GetValuesfromStructFromRedis func does get value from Redis
func GetBooksValuesFromRedis(cmd redis.Conn) error {

	// Set keys value
	booksKeys := "Book-1"

	// Get value from Redis
	valuesBooks, err := cmd.Do("GET", booksKeys)

	// Exception error handling
	if err != nil {
		panic(err.Error())
	} else if valuesBooks == nil {
		panic(err)
	} else {
		fmt.Printf("%s = %s\n", booksKeys, valuesBooks)
	}

	// Returned nil
	return nil

}
