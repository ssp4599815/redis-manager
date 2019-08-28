package rdb

import (
	"github.com/dongmx/rdb/nopdecoder"
	"strconv"
)

// Entry is info of a redis recored
type Entry struct {
	Key                string // key 名
	Bytes              uint64 // value 的大小
	Type               string // key 的类型
	NumOfElem          uint64 // value 中元素的个数
	LenOfLargestElem   uint64 // 最长元素的长度
	FieldOfLargestElem string // 最长元素的名称
}

// Decoder decode rdb file
type Decoder struct {
	Entries  chan *Entry
	m        MemProfiler
	tmpStore map[string]*Entry
	nopdecoder.NopDecoder
}

// NewDecoder new a rdb decoder
func NewDecoder() *Decoder {
	return &Decoder{
		Entries:  make(chan *Entry, 1024), // 初始化一个通道
		m:        MemProfiler{},
		tmpStore: map[string]*Entry{},
	}
}

// Set is called once for each string key
func (d *Decoder) Set(key, value []byte, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += d.m.SizeofString(value)
	bytes += d.m.TopLevelObjOverhead()
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "string",
		NumOfElem: d.m.ElemLen(value),
	}
	d.Entries <- e
}

// StartHash is called at the beginning of a hash.
// Hset will be called exactly length times before EndHash.
func (d *Decoder) StartHash(key []byte, length, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.HashtableOverhead(uint64(length))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "hash",
		NumOfElem: uint64(length),
	}
	d.tmpStore[keyStr] = e
}

// Hset is called once for each field=value pair in a hash.
func (d *Decoder) Hset(key, field, value []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(field) + d.m.ElemLen(value)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(field)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += d.m.SizeofString(field)
	e.Bytes += d.m.SizeofString(value)
	e.Bytes += d.m.HashtableEntryOverhead()
	e.Bytes += 2 * d.m.RobjOverhead()
	d.tmpStore[keyStr] = e
}

// EndHash is called when there are no more fields in a hash.
func (d *Decoder) EndHash(key []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	d.Entries <- e
	delete(d.tmpStore, keyStr)
}

// StartSet is called at the beginning of a set.
// Sadd will be called exactly cardinality times before EndSet.
func (d *Decoder) StartSet(key []byte, cardinality, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.HashtableOverhead(uint64(cardinality))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "set",
		NumOfElem: uint64(cardinality),
	}
	d.tmpStore[keyStr] = e
}

// Sadd is called once for each member of a set.
func (d *Decoder) Sadd(key, member []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += d.m.SizeofString(member)
	e.Bytes += d.m.HashtableEntryOverhead()
	e.Bytes += d.m.RobjOverhead()
	d.tmpStore[keyStr] = e
}

// EndSet is called when there are no more fields in a set.
// Same as EndHash
func (d *Decoder) EndSet(key []byte) {
	d.EndHash(key)
}

// StartList is called at the beginning of a list.
// Rpush will be called exactly length times before EndList.
// If length of the list is not known, then length is -1
func (d *Decoder) StartList(key []byte, length, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.LinkedListEntryOverhead() * uint64(length)
	bytes += d.m.LinkedlistOverhead()
	bytes += d.m.RobjOverhead() * uint64(length)
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "list",
		NumOfElem: uint64(length),
	}
	d.tmpStore[keyStr] = e
}

// Rpush is called once for each value in a list.
func (d *Decoder) Rpush(key, value []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(value)
	if _, err := strconv.ParseInt(string(value), 10, 32); err == nil {
		e.Bytes += 4
	} else {
		e.Bytes += d.m.SizeofString(value)
	}
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(value)
		e.LenOfLargestElem = lenOfElem
	}
	d.tmpStore[keyStr] = e
}

// EndList is called when there are no more values in a list.
func (d *Decoder) EndList(key []byte) {
	d.EndHash(key)
}

// StartZSet is called at the beginning of a sorted set.
// Zadd will be called exactly cardinality times before EndZSet.
func (d *Decoder) StartZSet(key []byte, cardinality, expiry int64) {
	keyStr := string(key)
	bytes := d.m.SizeofString(key)
	bytes += 2 * d.m.RobjOverhead()
	bytes += d.m.TopLevelObjOverhead()
	bytes += d.m.KeyExpiryOverhead(expiry)
	bytes += d.m.SkiplistOverhead(uint64(cardinality))
	e := &Entry{
		Key:       keyStr,
		Bytes:     bytes,
		Type:      "sortedset",
		NumOfElem: uint64(cardinality),
	}
	d.tmpStore[keyStr] = e
}

// Zadd is called once for each member of a sorted set.
func (d *Decoder) Zadd(key []byte, score float64, member []byte) {
	keyStr := string(key)
	e := d.tmpStore[keyStr]
	lenOfElem := d.m.ElemLen(member)
	if lenOfElem > e.LenOfLargestElem {
		e.FieldOfLargestElem = string(member)
		e.LenOfLargestElem = lenOfElem
	}
	e.Bytes += 8 // sizeof(score)
	e.Bytes += d.m.SizeofString(member)
	e.Bytes += 2 * d.m.RobjOverhead()
	e.Bytes += d.m.SkiplistEntryOverhead()
	d.tmpStore[keyStr] = e
}

// EndZSet is called when there are no more members in a sorted set.
func (d *Decoder) EndZSet(key []byte) {
	d.EndHash(key)
}

// EndRDB is called when parsing of the RDB file is complete.
func (d *Decoder) EndRDB() {
	close(d.Entries)
}
