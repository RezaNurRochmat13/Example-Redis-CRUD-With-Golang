package main

import (
	"github.com/gin-gonic/gin"
	"go-example-redis/controller"
)

// MAIN APPS
func main() {
	BaseRouting()
}

// BASE ROUTING APPS
func BaseRouting() {
	routers := gin.Default()

	v1 := routers.Group("/v1/api/")
	{
		v1.GET("books/:UUIDBooks", controller.GetDetailBooksFromRedis)
		v1.POST("books", controller.CreateBooksToDBAndRedis)
		v1.PUT("books/:UUIDBooks", controller.UpdateBooksInRedisAndDB)
		v1.DELETE("books/:UUIDBooks", controller.DeleteBooksFromRedisAndDB)
	}

	routers.Run(":8400")
}
