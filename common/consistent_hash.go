package common

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return "ConsistentError: " + e.s
}

// 定义错误类型
func ConsistentError(text string) error {
	return &errorString{text}
}

// 定义环类型
type Circle []uint32

func (c Circle) Len() int {
	return len(c)
}

func (c Circle) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c Circle) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Hash func(date []byte) uint32

type Consistent struct {
	hash         Hash              // 产生uint32类型的函数
	circle       Circle            // 环
	virtualNodes int               // 虚拟节点个数
	virtualMap   map[uint32]string // 点到主机的映射
	members      map[string]bool   // 主机列表
	sync.RWMutex
}

func NewConsisten() *Consistent {
	return &Consistent{
		hash:         crc32.ChecksumIEEE,
		circle:       Circle{},
		virtualNodes: 150,
		virtualMap:   make(map[uint32]string),
		members:      make(map[string]bool),
	}
}

// generate a string key for an element with an index
func (c *Consistent) eltKey(key string, idx int) string {
	return key + "|" + strconv.Itoa(idx)
}

func (c *Consistent) updateCricle() {
	c.circle = Circle{}
	for k := range c.virtualMap {
		c.circle = append(c.circle, k)
	}
	sort.Sort(c.circle)
}

func (c *Consistent) Members() []string {
	c.RLock()
	defer c.RUnlock()

	m := make([]string, len(c.members))

	var i = 0
	for k := range c.members {
		m[i] = k
		i++
	}

	return m
}

func (c *Consistent) Get(key string) string {
	hashKey := c.hash([]byte(key))
	c.RLock()
	defer c.RUnlock()

	i := c.search(hashKey)
	//fmt.Println("i", i)

	return c.virtualMap[c.circle[i]]
}

// search nearly vnode around key
// sort.Search uses binary search to find key
// every vnode cover its self and clockwise area
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.circle[x] >= key
	}

	i := sort.Search(len(c.circle), f)
	i = i - 1
	if i < 0 {
		i = len(c.circle) - 1
	}
	return i
}

// this function is beautiful
func (c *Consistent) ForceSet(keys ...string) {
	mems := c.Members()
	for _, elt := range mems {
		var found = false

	FOUNDLOOP:
		for _, k := range keys {
			if k == elt {
				found = true
				break FOUNDLOOP
			}
		}
		if !found {
			c.Remove(elt)
		}
	}

	for _, k := range keys {
		c.RLock()
		_, ok := c.members[k]
		c.RUnlock()

		if !ok {
			c.Add(k)
		}
	}
}

func (c *Consistent) Add(elt string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.members[elt]; ok {
		return
	}

	c.members[elt] = true

	for idx := 0; idx < c.virtualNodes; idx++ {
		c.virtualMap[c.hash([]byte(c.eltKey(elt, idx)))] = elt
	}

	c.updateCricle()
}

func (c *Consistent) Remove(elt string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.members[elt]; !ok {
		return
	}

	delete(c.members, elt)

	for idx := 0; idx < c.virtualNodes; idx++ {
		delete(c.virtualMap, c.hash([]byte(c.eltKey(elt, idx))))
	}

	c.updateCricle()
}
