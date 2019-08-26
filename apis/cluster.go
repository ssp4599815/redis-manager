package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/ssp4599815/redis-manager/utils"
	"net/http"
)

func Nodes(c *gin.Context) {
	var info utils.Info
	c.Bind(&info)
	nodes := utils.ParseClusterNodes(info)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "获取nodes成功",
		"data":    nodes,
	})
}
