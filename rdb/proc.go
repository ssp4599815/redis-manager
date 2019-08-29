package rdb

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dongmx/rdb"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	byteLF = byte('\n') // 换行符
)

var (
	bytesSpace       = []byte(" ")
	psyncFullSyncCmd = []byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n") // 执行全量复制的命令
)

// Instance is the struct for instance node
type Instance struct {
	Addr   string
	Target string

	tconn net.Conn

	lock sync.RWMutex
	conn net.Conn
	br   *bufio.Reader
	bw   *bufio.Writer

	barrierC chan struct{}
	wg       *sync.WaitGroup

	offset   int64
	masterID string
}

func (inst *Instance) Sync() {
	defer inst.wg.Done()
	fmt.Printf("tring to sync with remote instance %s\n", inst.Addr)
	// <-inst.barrierC // TODO 啥意思

	for {
		err := inst.sync()
		if err != nil {
			fmt.Printf("fail to syncing redis data due %s\n", err)
		}
		time.Sleep(time.Second * 30)
	}
}

func (inst *Instance) sync() (err error) {
	fmt.Printf("starting to sync with remote instance %s\n", inst.Addr)
	defer inst.Close()

	atomic.StoreInt64(&inst.offset, 0)

	conn, err := net.Dial("tcp", inst.Addr) // 和远端 Redis 建立连接
	if err != nil {
		return err
	}
	inst.conn = conn
	inst.bw = bufio.NewWriter(conn)
	inst.br = bufio.NewReader(conn)

	// 1. barrier run syncRDB
	// 1.1 send psync ? -1 表示需要全量复制
	fmt.Printf("start to sync rdb of %s\n", inst.Addr)
	_ = writeAll(psyncFullSyncCmd, inst.bw) // 向 conn 中写入数据
	_ = inst.bw.Flush()                     // 刷新数据到磁盘
	var data []byte
	data, err = inst.br.ReadBytes(byteLF) // 读取文件直到换行符
	if err != nil {
		return
	}

	fmt.Printf("parse paync reply of %s\n", inst.Addr)
	err = inst.parsePSyncReply(data) // data 为服务端响应的数据的第一行 （+FULLRESYNC d3c009753b7dc26efbf1b6b04e63a93bcfbfb4b1 0）
	if err != nil {
		return
	}

	// because rdb was transformed by RESP Bulk String, we need ignore first line
	// 忽略第一行的解析数据
	/*
		执行命令： PSYNC ? -1
		+FULLRESYNC d3c009753b7dc26efbf1b6b04e63a93bcfbfb4b1 0
		$206
		REDIS0009�	redis-ver5.0.0�
		...
	*/
	for {
		data, err = inst.br.ReadBytes(byteLF) // 读取到换行符
		if err != nil {
			return
		}
		// func Quote(s string) string
		// 返回字符串s在go语法下的双引号字面值表示，控制字符、不可打印字符会进行转义。（如\t，\n，\xFF，\u0100）
		fmt.Printf("read new line add %s with %s\n", inst.Addr, strconv.Quote(string(data)))
		if len(data) > 0 && data[0] == byte('$') { // 一般的 结束符为 $206
			fmt.Println("读取完PSYNC返回的信息，跳出循环")
			break
		}
	}

	// read full rdb
	//err = inst.syncRDB(inst.Target)

	//for {
	//	data, err = inst.br.ReadBytes(byteLF) // 读取到换行符
	//	fmt.Printf(string(data))
	//	time.Sleep(time.Second)
	//}
	decoder := NewDecoder()
	err = rdb.Decode(inst.br, decoder)
	cnt := NewCounter()
	cnt.Count(decoder.Entries)
	saa := GetData(cnt)
	jsonBytes, _ := json.MarshalIndent(saa, "", "    ")
	fmt.Println(string(jsonBytes))

	return
}

func (inst *Instance) parsePSyncReply(data []byte) error {
	fmt.Printf("receive paync data %s reply as %s\n", inst.Addr, strconv.Quote(string(data)))
	splited := bytes.Split(data, bytesSpace)           // 根据空格进行分割
	runidBs := string(splited[1])                      // 运行的当前master id
	offsetBs := string(splited[2][:len(splited[2])-2]) // 偏移量

	offset, err := strconv.ParseInt(offsetBs, 10, 64)
	if err != nil {
		return err
	}
	inst.offset = offset
	inst.masterID = runidBs
	return nil
}

// Close the up and down steam
func (inst *Instance) Close() {
	if inst.conn != nil {
		inst.conn.Close()
	}
	if inst.tconn != nil {
		inst.tconn.Close()
	}
	return
}

func (inst *Instance) syncRDB(addr string) (err error) {
	fmt.Printf("start syncing rdb for %s", inst.Addr)
	return nil
}

func writeAll(buf []byte, w io.Writer) error {
	left := len(buf) // buf 中剩余的数据
	for left != 0 {
		size, err := w.Write(buf[len(buf)-left:])
		if err != nil {
			return err
		}
		left -= size
	}
	return nil
}
