package apis

import (
	"testing"
)

func TestParseClusterNodes(t *testing.T) {
	ParseClusterNodes("10.211.55.12", "8001", "")
}