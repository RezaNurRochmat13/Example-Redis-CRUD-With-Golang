package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-redis/model"
)

func main() {
	// Initialized database pooling
	initializePools := newPools()
	establishedConn := initializePools.Get()
	// Closing database connection
	defer establishedConn.Close()

	// Check ping Redis
	err := pingConnection(establishedConn)
	if err != nil {
		panic(err)
	}

	// Set string values to Redis
	err = SetValuesToRedis(establishedConn)

	if err != nil {
		panic(err)
	}

	// Get string values from Redis
	err = GetValuesFromRedis(establishedConn)

	if err != nil {
		panic(err)
	}

	// Set values from struct in Redis
	err = setValuesFromStruct(establishedConn)

	if err != nil {
		panic(err)
	}

	// Get values from struct in Redis
	err = getValuesFromStruct(establishedConn)

	if err != nil {
		panic(err)
	}

}

// Configure database pooling
func newPools() *redis.Pool {
	return &redis.Pool{
		// Set maximum idle connection
		MaxIdle: 80,
		// Set max number of connection
		MaxActive: 12000,
		// Configure connection
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")

			if err != nil {
				panic(err.Error())
			}

			return c, err
		},
	}
}

// Ping check connection to Redis
func pingConnection(cmd redis.Conn) error {
	// Send PING command to Redis
	pong, err := cmd.Do("PING")

	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple string"
	// use redis.String to convert to interface type to string
	strings, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING response is : %s\n", strings)
	// Output PONG

	return nil
}

// Set values into Redis
func SetValuesToRedis(cmd redis.Conn) error {
	_, err := cmd.Do("SET", "Favorite Language", "Golang")

	if err != nil {
		return err
	}

	_, err = cmd.Do("SET", "Favorite Database", "MongoDB")

	if err != nil {
		return err
	}

	return nil
}

// Get values from Redis
func GetValuesFromRedis(cmd redis.Conn) error {
	// Simple GET example with String helper
	keys := "Favorite Language"
	strings, err := redis.String(cmd.Do("GET", keys))

	if err != nil {
		return err
	}

	fmt.Printf("%s = %s\n", keys, strings)

	keys2 := "Favorite Database"
	string2, err := redis.String(cmd.Do("GET", keys2))

	if err != nil {
		return nil
	}

	fmt.Printf("%s = %s\n", keys2, string2)

	return nil
}

// Set values from struct in Redis
func setValuesFromStruct(cmd redis.Conn) error {
	user := &model.Users{
		Username: "Rejak",
		MobileID: "12345678",
		Email: "rejak@gmail.com",
		Firstname: "Rejak",
		Lastname: "Rochmat",
	}

	serialized, err := json.Marshal(user)

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Serialized data : ", string(serialized))
	}

	_, err = cmd.Do("SET", "1", serialized)

	if err != nil {
		return err
	}

	return nil
}

// Get values from struct in Redis
func getValuesFromStruct(cmd redis.Conn) error {

	keys := "1"

	strings, err := cmd.Do("GET", keys)

	if err != nil {
		return err
	}

	fmt.Printf("%s = %s\n", keys, strings)

	return nil
}


