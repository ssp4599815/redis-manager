package apis

import (
	"container/heap"
	"fmt"
	"sort"
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
	heap.Init(h)
	p := &prefixHeap{}
	heap.Init(p)
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
		separators:         ";:,_-",
	}
}

// Count by variout dimensions
func (c *Counter) Count(in <-chan *Entry) {
	for e := range in {
		c.count(e)
	}
	// get largest prefixes

}

func (c *Counter) count(e *Entry) {
	fmt.Println(e)
	c.countLargestEntries(e, 500)
}

func (c *Counter) countLargestEntries(e *Entry, num int) {
	heap.Push(c.largestEntries, e)
	l := c.largestEntries.Len()
	if l > num {
		heap.Pop(c.largestEntries)
	}
}

//GetLargestEntries from heap ,num max is 500
func (c *Counter) GetLargestEntries(num int) []*Entry {
	res := []*Entry{}

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
	typeKey
	Bytes uint64
	Num   uint64
}

func (h prefixHeap) Len() int {
	return len(h)
}
func (h prefixHeap) Less(i, j int) bool {
	if h[i].Bytes < h[j].Bytes {
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
