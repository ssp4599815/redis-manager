package rdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePSyncReply(t *testing.T) {
	data := []byte("+fullsync 0123456789012345678901234567890123456789 7788\r\n")
	inst := &Instance{}
	err := inst.parsePSyncReply(data)
	assert.NoError(t, err)
}

func TestSync(t *testing.T) {
	inst := &Instance{
		Addr: "127.0.0.1:6379",
	}
	inst.Sync()
}
