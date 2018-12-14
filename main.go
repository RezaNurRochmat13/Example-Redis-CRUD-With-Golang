package main

import (
	"github.com/gin-gonic/gin"
	"go-example-redis/controller"
)

func main() {
	BaseRouting()

}

func BaseRouting() {
	routers := gin.Default()

	v1 := routers.Group("/v1/api/")
	{
		v1.GET("books/:UUIDBooks", controller.GetDetailBooksFromRedis)
		v1.POST("books", controller.CreateBooksToDBAndRedis)
		v1.PUT("books/:UUIDBooks", controller.UpdateBooksInRedisAndDB)
	}

	routers.Run(":8400")
}
