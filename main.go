package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/ssp4599815/redis-manager/apis"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/api/v1/detail", GetClusterInfos)

	return router
}

func main() {
	r := initRouter()
	r.Run(":8088")
}
