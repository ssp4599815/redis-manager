package apis

import (
	"fmt"
	"github.com/dongmx/rdb"
	"os"
)

func decode(decoder *Decoder, filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("open rdbfile err: ", err)
		return
	}
	err = rdb.Decode(f, decoder)
	if err != nil {
		fmt.Println("decode rdbfile err: ", err)
		return
	}
}

func Dump() {
	decoder := NewDecoder()
	go decode(decoder, "/Users/a1800101257/go/src/github.com/ssp4599815/redis-manager/apis/rdb.rdb")
	cnt := NewCounter()
	cnt.Count(decoder.Entries)
}
