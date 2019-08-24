package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

func parseClusterNodes() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       []string{
			"127.0.0.1:7000",
			"127.0.0.1:7001",
			"127.0.0.1:7002",
			"127.0.0.1:7003",
			"127.0.0.1:7004",
			"127.0.0.1:7005",
		},
		ReadTimeout: time.Second * 3,
		DialTimeout: 3 * time.Second,
	})

	response := client.ClusterNodes()
	fmt.Println(response)
	raw_lines := response.Val()
	for _, line := range strings.Split(raw_lines, "\n") {
		if len(line) == 0 { // 跳过空行
			continue
		}
		parseNodeLine(line)
	}

}
func parseNodeLine(line string) {
	fmt.Println(line)
}
func main() {
	parseClusterNodes()
}
