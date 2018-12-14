package config

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// DatabaseConn function does initialize database connection
func DatabaseConn() *gorm.DB {
	db, err := gorm.Open("mysql", "reja:reja123@tcp(192.168.203.86:3306)/redis-database?charset=utf8")

	if err != nil {
		panic(err.Error())
	}

	return db
}

func NewConnectionPools() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
