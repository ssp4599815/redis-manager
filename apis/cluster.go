package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strings"
)

func GetClusterInfos(c *gin.Context) {

	host := "192.168.16.56"
	port := "8001"
	password := ""
	nodes := ParseClusterNodes(host, port, password)
	fmt.Println(nodes)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "get cluster info successful",
		"data":    nodes,
	})
}

func ParseClusterNodes(host, port, password string) map[string]node {
	addr := strings.Join([]string{host, port}, ":")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr},
		Password: password,
	})
	response := client.ClusterNodes()
	rawLines := response.Val()

	nodes := make(map[string]node)
	for _, line := range strings.Split(rawLines, "\n") {
		if len(line) == 0 { // 跳过空行
			continue
		}
		nl := parseNodeLine(line)
		nodes[nl.addr] = nl
	}
	return nodes
}

type node struct {
	nodeId    string
	addr      string
	flags     string
	masterId  string
	connected bool
}

func parseNodeLine(line string) node {
	var n node
	lineItmes := strings.Split(line, " ")
	n.nodeId = lineItmes[0]
	n.addr = parseAddr(lineItmes[1])
	n.flags = parseFlags(lineItmes[2])
	n.masterId = parseMasterId(lineItmes[3])
	n.connected = parseConnected(lineItmes[7])
	return n
}
func parseAddr(addr string) string {
	addr = strings.Split(addr, "@")[0]
	return addr
}

func parseFlags(flags string) string {
	if strings.HasPrefix(flags, "myself") {
		flags = strings.Split(flags, ",")[1]
	}
	return flags
}
func parseMasterId(masterId string) string {
	if masterId == "-" {
		return ""
	} else {
		return masterId
	}
}

func parseConnected(connected string) bool {
	if connected == "connected" {
		return true
	} else {
		return false
	}
}
