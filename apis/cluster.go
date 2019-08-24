package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strings"
)

func GetClusterInfos(c *gin.Context) {

	host := "127.0.0.1"
	port := "7000"
	password := ""
	ParseClusterNodes(host, port, password)
	fmt.Println(123)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get cluster info successful",
		"data":    "ok",
	})
}

func ParseClusterNodes(host, port, password string) {
	addr := strings.Join([]string{host, port}, ":")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr},
		Password: password,
	})
	fmt.Println(client.ClusterNodes())
	response := client.ClusterNodes()
	raw_lines := response.Val()
	for _, line := range strings.Split(raw_lines, "\n") {
		if len(line) == 0 { // 跳过空行
			continue
		}
		fmt.Println(line)
		parseNodeLine(line)
	}

}
func parseNodeLine(line string) {
	fmt.Println(line)
}
