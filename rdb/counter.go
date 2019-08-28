package rdb

import (
	"container/heap"
	"sort"
	"strconv"
	"strings"
)

// Counter for redis memory userage
type Counter struct {
	largestEntries     *entryHeap
	largestKeyPrefixed *prefixHeap
	lengthLevel0       uint64
	lengthLevel1       uint64
	lengthLevel2       uint64
	lengthLevel3       uint64
	lengthLevel4       uint64
	lengthLevelBytes   map[typeKey]uint64
	lengthLevelNum     map[typeKey]uint64
	keyPrefixBytes     map[typeKey]uint64
	keyPrefixNum       map[typeKey]uint64
	separators         string
	typeBytes          map[string]uint64
	typeNum            map[string]uint64
}

// NewCounter reture a pointer of Counter
func NewCounter() *Counter {
	h := &entryHeap{}
	heap.Init(h) // 初始化 堆
	p := &prefixHeap{}
	heap.Init(p) // 初始化 堆
	return &Counter{
		largestEntries:     h,
		largestKeyPrefixed: p,
		lengthLevel0:       100,
		lengthLevel1:       1000,
		lengthLevel2:       10000,
		lengthLevel3:       100000,
		lengthLevel4:       1000000,
		lengthLevelBytes:   map[typeKey]uint64{},
		lengthLevelNum:     map[typeKey]uint64{},
		keyPrefixBytes:     map[typeKey]uint64{},
		keyPrefixNum:       map[typeKey]uint64{},
		typeBytes:          map[string]uint64{},
		typeNum:            map[string]uint64{},
		separators:         ";:,_-", // 分隔符，用来统计 不同类型的 key 前缀
	}
}

// Count by variout dimensions
func (c *Counter) Count(in <-chan *Entry) { // Entry 是解析出来的一个 Redis 信息
	for e := range in { // 遍历通道来获取数据
		c.count(e)
	}
	// get largest prefixes
	c.calcuLargestKeyPrefix(1000)
}

func (c *Counter) count(e *Entry) {
	// 统计 topK 根据 key 的大小。
	c.countLargestEntries(e, 500)
	// 统计key的前缀个数
	c.countByKeyPrefix(e)
	// 统计不同类型key 的 个数和大小
	c.countByType(e)
	// 根据长度进行排序
	c.countByLength(e)
}

func (c *Counter) countByType(e *Entry) {
	c.typeNum[e.Type]++            // 统计不同类型的个数
	c.typeBytes[e.Type] += e.Bytes // 统计不同类型的大小
}

func (c *Counter) countByLength(e *Entry) {
	key := typeKey{
		Type: e.Type,
		// func FormatInt(i int64, base int) string
		// 返回i的base进制的字符串表示。base 必须在2到36之间，结果中会使用小写字母'a'到'z'表示大于10的数字。
		Key: strconv.FormatUint(c.lengthLevel0, 10),
	}
	add := func(c *Counter, key typeKey, e *Entry) {
		c.lengthLevelBytes[key] += e.Bytes
		c.lengthLevelNum[key]++
	}
	// must lengthLevel4 > lengthLevel3 > lengthLevel2 ...
	if e.NumOfElem > c.lengthLevel4 {
		key.Key = strconv.FormatUint(c.lengthLevel4, 10)
		add(c, key, e)
	} else if e.NumOfElem > c.lengthLevel3 {
		key.Key = strconv.FormatUint(c.lengthLevel3, 10)
		add(c, key, e)
	} else if e.NumOfElem > c.lengthLevel2 {
		key.Key = strconv.FormatUint(c.lengthLevel2, 10)
		add(c, key, e)
	} else if e.NumOfElem > c.lengthLevel1 {
		key.Key = strconv.FormatUint(c.lengthLevel1, 10)
		add(c, key, e)
	} else if e.NumOfElem > c.lengthLevel0 {
		key.Key = strconv.FormatUint(c.lengthLevel0, 10)
		add(c, key, e)
	}
}

func (c *Counter) countByKeyPrefix(e *Entry) {
	// reset all numbers to 0
	k := strings.Map(func(r rune) rune {
		if r >= 48 && r <= 57 { //48 == "0" 57 == "9"
			return '0'
		}
		return r
	}, e.Key)
	// 统计不同类型的前缀
	prefixes := getPrefixes(k, c.separators)
	key := typeKey{
		Type: e.Type,
	}
	for _, prefix := range prefixes {
		if len(prefix) == 0 {
			continue
		}
		key.Key = prefix
		c.keyPrefixBytes[key] += e.Bytes // 统计 key prefix 的大小
		c.keyPrefixNum[key]++            // 统计 key prefix 的个数
	}
}

func (c *Counter) countLargestEntries(e *Entry, num int) {
	heap.Push(c.largestEntries, e) // 向堆中插入一个元素
	l := c.largestEntries.Len()
	if l > num {
		heap.Pop(c.largestEntries) // 如果达到 topK 的长度，就 弹出一个元素
	}
}

// 计算key prefix 的最大值，使用 堆排序
func (c *Counter) calcuLargestKeyPrefix(num int) {
	// keyPrefixBytes 为各种 key prefix 的 大小
	for key := range c.keyPrefixBytes { //  遍历的是一个 结构体的 key
		k := &PrefixEntry{} // 创建一个 prefix 堆用来排序
		k.Type = key.Type
		k.Key = key.Key
		k.Bytes = c.keyPrefixNum[key] // key prefix 的大小
		k.Num = c.keyPrefixNum[key]   // key prefix 的个数
		delete(c.keyPrefixBytes, key)
		delete(c.keyPrefixNum, key)
		heap.Push(c.largestKeyPrefixed, k) // 向堆中 push 一个元素
		l := c.largestKeyPrefixed.Len()    // 超过限额就pop 一个元素
		if l > num {
			heap.Pop(c.largestKeyPrefixed)
		}
	}
}

// GetLargestKeyPrefixes form heap
func (c *Counter) GetLargestKeyPrefixes() []*PrefixEntry {
	var res []*PrefixEntry

	// get a copy of c.largestKeyPrefixs
	for i := 0; i < c.largestKeyPrefixed.Len(); i++ {
		entries := *c.largestKeyPrefixed
		res = append(res, entries[i])
	}
	sort.Sort(sort.Reverse(prefixHeap(res)))
	return res
}

func (c *Counter) GetLenLevelCount() []*PrefixEntry {
	var res []*PrefixEntry

	// get a copy of lengthLevelBytes and lengthLevelNum
	for key := range c.lengthLevelBytes {
		entry := &PrefixEntry{}
		entry.Type = key.Type
		entry.Key = key.Key
		entry.Bytes = c.lengthLevelBytes[key]
		entry.Num = c.lengthLevelNum[key]
		res = append(res, entry)
	}
	return res
}

//GetLargestEntries from heap ,num max is 500
func (c *Counter) GetLargestEntries(num int) []*Entry {
	var res []*Entry

	//  get a copy of c.largestEntries
	for i := 0; i < c.largestEntries.Len(); i++ {
		entries := *c.largestEntries
		res = append(res, entries[i])
	}
	sort.Sort(sort.Reverse(entryHeap(res)))
	if num < len(res) {
		res = res[:num]
	}
	return res
}

type typeKey struct {
	Type string
	Key  string
}

type entryHeap []*Entry

func (h entryHeap) Len() int {
	return len(h)
}
func (h entryHeap) Less(i, j int) bool {
	return h[i].Bytes < h[j].Bytes
}
func (h entryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *entryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h *entryHeap) Push(e interface{}) {
	*h = append(*h, e.(*Entry))
}

type prefixHeap []*PrefixEntry

type PrefixEntry struct {
	// 用作统计前缀的 结构体
	typeKey      // key 的类型和值
	Bytes uint64 // 前缀的大小
	Num   uint64 // 前缀的个数
}

func (h prefixHeap) Len() int {
	return len(h)
}

// 小顶堆
func (h prefixHeap) Less(i, j int) bool {
	if h[i].Bytes < h[j].Bytes { // 用于比较的 两个元素
		return true
	} else if h[i].Bytes == h[j].Bytes {
		if h[i].Num < h[j].Num {
			return true
		} else if h[i].Num == h[j].Num {
			if h[i].Key > h[j].Key {
				return true
			}
		}
	}
	return false
}
func (h prefixHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *prefixHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h *prefixHeap) Push(k interface{}) {
	*h = append(*h, k.(*PrefixEntry))
}

// 获取前缀
func getPrefixes(s, sep string) []string {
	var res []string

	// key中是否有 sep（一个分隔符）,有就取出索引，没有就返回-1
	sepIdx := strings.IndexAny(s, sep)

	if sepIdx < 0 { // 没有找到 分隔符
		res = append(res, s)
	}
	// 如果找到了 分隔符的索引位置
	for sepIdx > -1 {
		r := s[:sepIdx+1] // 取出分隔符前面的 key，包括分隔符
		if len(res) > 0 { // 如果 res 不为空，也就是说 里面有了各种key 的前缀
			// res[len(res)-1 ] 是 获取前面的key前缀
			// s[:sepIdx+1] 是 获取 新key前缀到分隔符的 字符串
			r = res[len(res)-1 ] + s[:sepIdx+1] // 组合成新的 key 前缀
		}
		res = append(res, r)
		s = s[sepIdx+1:]                  // 重置key，根据分隔符 把分隔符前面的字符串去掉，生成一个新的key
		sepIdx = strings.IndexAny(s, sep) // 重置key，重新获取 seq 的索引位置
	}
	// 去掉 带分隔符的后缀
	// Trim all suffix of  separators
	for i := range res { // 获取每一个 包含前缀的 key
		for hasAnySuffix(res[i], sep) { // 是否包含 带分隔符的后缀
			res[i] = res[i][:len(res[i])-1]
		}
	}
	res = removeDuplicatesUnordered(res)
	return res
}

// 移除 重复的 无序的元素
func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// 元素去重
	// create a map of all unique elements
	for v := range elements {
		encountered[elements[v]] = true
	}

	var result []string
	for key := range encountered {
		result = append(result, key)
	}

	return result
}

// 后缀判断
func hasAnySuffix(s, suffix string) bool {
	for _, c := range suffix {
		if strings.HasSuffix(s, string(c)) {
			return true
		}
	}
	return false
}
