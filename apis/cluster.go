package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strings"
)

func GetClusterNodes(c *gin.Context) {
	fmt.Println(c.PostForm("host"))
	fmt.Println(c.Request.PostForm)
	fmt.Println(c.PostFormArray("host"))
	fmt.Println(c.PostFormMap("host"))
	fmt.Println(c.Request.Header)
	buf := make([]byte, 1024)
	n,_ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[0:n]))

	fmt.Println(c.JSON)



	host := "10.211.55.12"
	port := "8001"
	password := ""

	nodes := ParseClusterNodes(host, port, password)
	_, err := c.Writer.WriteString(nodes)
	if err != nil {
		fmt.Println("response write failed,err: ", err)
	}
}

func ParseClusterNodes(host, port, password string) string {
	addr := strings.Join([]string{host, port}, ":")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr},
		Password: password,
	})
	response := client.ClusterNodes()
	rawLines := response.Val()
	nodes := make([]map[string]interface{}, 0)

	for _, line := range strings.Split(rawLines, "\n") {
		if len(line) == 0 { // 跳过空行
			continue
		}
		nl := parseNodeLine(line)
		nodes = append(nodes, nl)
	}
	data, _ := json.Marshal(nodes)
	//fmt.Println(string(data))
	return string(data)
}

func parseNodeLine(line string) map[string]interface{} {
	node := make(map[string]interface{})
	lineItmes := strings.Split(line, " ")
	node["NodeId"] = lineItmes[0]
	node["Addr"] = parseAddr(lineItmes[1])
	node["Flags"] = parseFlags(lineItmes[2])
	node["MasterId"] = parseMasterId(lineItmes[3])
	node["Connected"] = parseConnected(lineItmes[7])
	return node
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
