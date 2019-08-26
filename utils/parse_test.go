package utils

import "testing"

func TestParseClusterNodes(t *testing.T) {
	info := Info{"loan", "10.211.55.12", "8001", ""}
	ParseClusterNodes(info)
}
