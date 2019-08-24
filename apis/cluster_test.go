package apis

import (
	"testing"
)

func TestParseClusterNodes(t *testing.T) {
	ParseClusterNodes("192.168.16.56", "8001", "")
}