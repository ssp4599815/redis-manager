package rdb

import (
	"testing"
)

func TestCounter_GetLargestEntries(t *testing.T) {
	//e := &Entry{
	//	Key: "RELATIONSFOLLOWERIDS6420000664",
	//	Bytes: 1,
	//	Type: "sortedset",
	//	NumOfElem: 1,
	//	LenOfLargestElem:1,
	//	FieldOfLargestElem: "test",
	//}
	c := NewCounter()
	decoder := NewDecoder()
	c.Count(decoder.Entries)
}
