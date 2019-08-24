package apis

import (
	"fmt"
	"testing"
)

func TestParseClusterNodes(t *testing.T) {
	ParseClusterNodes("127.0.0.1", "7000", "")
	fmt.Println("123")
}