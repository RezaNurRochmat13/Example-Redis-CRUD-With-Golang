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
		v1.POST("books", controller.SetBooksValuesToDBAndRedis)
	}

	routers.Run(":8400")
}
