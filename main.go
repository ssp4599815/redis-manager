package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/ssp4599815/redis-manager/apis"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/v1/nodes", GetClusterNodes)
	return router
}

func main() {
	r := initRouter()
	r.Run(":8089")
}
