package rdb

import (
	"encoding/json"
	"fmt"
	"github.com/dongmx/rdb"
	"os"
	"path"
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
	baseDir, _ := os.Getwd()
	go decode(decoder, path.Join(baseDir, "rdb.rdb"))
	cnt := NewCounter()
	cnt.Count(decoder.Entries)
	data := GetData(cnt)
	jsonBytes, _ := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(jsonBytes))
}

func GetData(cnt *Counter) map[string]interface{} {
	data := make(map[string]interface{})
	data["LargestKeys"] = cnt.GetLargestEntries(100)

	largestKeyPrefixesByType := map[string][]*PrefixEntry{}
	for _, entry := range cnt.GetLargestKeyPrefixes() {
		// if mem usage is less than 1M, and the list is long enough, then it's unnecessary to add it.
		if entry.Bytes < 1000*1000 && len(largestKeyPrefixesByType[entry.Type]) > 50 {
			continue
		}
		largestKeyPrefixesByType[entry.Type] = append(largestKeyPrefixesByType[entry.Type], entry)
	}
	data["LargestKeyPrefixes"] = largestKeyPrefixesByType

	data["TypeBytes"] = cnt.typeBytes
	data["TypeNum"] = cnt.typeNum
	totalNum := uint64(0)
	for _, v := range cnt.typeNum {
		totalNum += v
	}

	totalBytes := uint64(0)
	for _, v := range cnt.typeBytes {
		totalBytes += v
	}

	data["TotalNum"] = totalNum
	data["TotalBytes"] = totalBytes

	lenLevelCount := map[string][]*PrefixEntry{}
	for _, entry := range cnt.GetLenLevelCount() {
		lenLevelCount[entry.Type] = append(lenLevelCount[entry.Type], entry)
	}
	data["LenLevelCount"] = lenLevelCount
	return data
}
