package consistenthash

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strconv"
	"sync"
)

var randBuf [16]byte

type Hash func(data []byte) uint32

type Node interface {
	Name() string
}

type ringNode struct {
	rid  uint32
	node Node
}

type CHash struct {
	sync.Mutex
	hash     Hash
	replicas int
	ring     []ringNode
	ringIds  map[uint32]struct{} // for make sure node id in ring is unique
}

func New(replicas int, fn Hash) *CHash {
	ch := &CHash{
		replicas: replicas,
		hash:     fn,
		ringIds:  make(map[uint32]struct{}),
	}
	if ch.hash == nil {
		ch.hash = crc32.ChecksumIEEE
	}
	return ch
}

//every node can has specific replicas, if not set, use ch.replicas as default
func (ch *CHash) AddNodes(nodes ...Node) {
	if len(nodes) == 0 {
		return
	}
	ch.Lock()
	defer ch.Unlock()

	var node_key string
	var rid uint32
	var newRing []ringNode
	for _, node := range nodes {
		rep := ch.replicas
		if r, ok := node.(interface {
			Replicas() int
		}); ok {
			rep = r.Replicas() //use node specific replicas
		}

		for i := 0; i < rep; i++ {
			node_key = node.Name() + "-" + strconv.Itoa(i)
		tryAgain:
			rid = ch.hash([]byte(node_key))
			if _, ok := ch.ringIds[rid]; ok {
				//exist,
				fmt.Printf("exsit, rid:%d, node_key:%s\n", rid, node_key)
				n, err := rand.Read(randBuf[:len(randBuf)-1])
				if err != nil {
					panic(err)
				}
				node_key = node_key + "-" + string(randBuf[:n])
				goto tryAgain
			}
			ch.ringIds[rid] = struct{}{}
			newRing = append(newRing, ringNode{rid: rid, node: node})
		}
	}
	newRing = append(newRing, ch.ring...)
	sort.Slice(newRing, func(i, j int) bool { return newRing[i].rid < newRing[j].rid })
	ch.ring = newRing //update
	//ch.ShowRing()
}

func (ch *CHash) GetNodeByKey(key uint32) Node {
	ring := ch.ring
	if len(ring) == 0 {
		return nil
	}

	idx := sort.Search(len(ring), func(i int) bool { return ring[i].rid >= key })
	if idx == len(ring) {
		idx = 0
	}

	return ring[idx].node
}

//if node is out of service, del the node from ring
func (ch *CHash) DelNode(delNode Node) {
	ch.Lock()
	defer ch.Unlock()

	newRing := []ringNode{}
	for _, rnode := range ch.ring {
		if rnode.node != delNode {
			newRing = append(newRing, rnode)
		}
	}
	ch.ring = newRing
}

func (ch *CHash) ShowRing() {
	for i, rnode := range ch.ring {
		fmt.Printf("i:%d, rid:%d, node.Name:%s\n", i, rnode.rid, rnode.node.Name())
	}
}

func (ch *CHash) GetRingLen() int {
	return len(ch.ring)
}

func (ch *CHash) GetRingNodeId(idx int) uint32 {
	if len(ch.ring) <= idx {
		return 0
	}
	return ch.ring[idx].rid
}
