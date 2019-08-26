package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ssp4599815/redis-manager/apis"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	// api路由分组
	api := router.Group("/api/v1")
	{
		api.POST("nodes", apis.Nodes)
	}
	return router
}
