package utils

import (
	"github.com/go-redis/redis"
	"strings"
)

type Info struct {
	Line     string `json:"line" form: "line"`
	Host     string `json:"host" form: "name"`
	Port     string `json:"port" form: "port"`
	Password string `json:"password,omitempty" form: "password"` // omitempty 可以忽略空值
}

func ParseClusterNodes(info Info) []Node {
	addr := strings.Join([]string{info.Host, info.Port}, ":")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr},
		Password: info.Password,
	})
	response := client.ClusterNodes()
	rawLines := response.Val()
	nodes := make([]Node, 0)

	for _, line := range strings.Split(rawLines, "\n") {
		if len(line) == 0 { // 跳过空行
			continue
		}
		nl := parseNodeLine(line)
		nodes = append(nodes, nl)
	}
	return nodes
}

type Node struct {
	NodeId    string
	Addr      string
	Flags     string
	MasterId  string
	Connected string
}

func parseNodeLine(line string) Node {
	lineItmes := strings.Split(line, " ")
	var node Node
	node.NodeId = lineItmes[0]
	node.Addr = parseAddr(lineItmes[1])
	node.Flags = parseFlags(lineItmes[2])
	node.MasterId = parseMasterId(lineItmes[3])
	node.Connected = parseConnected(lineItmes[7])
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
func parseConnected(connected string) string {
	if connected == "connected" {
		return "true"
	} else {
		return "false"
	}
}
