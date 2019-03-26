package consistenthash

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type serviceNode struct {
	name string
}

func (n *serviceNode) Name() string {
	return n.name
}

var myrand = rand.New(rand.NewSource(time.Now().Unix()))

func TestConsistency(t *testing.T) {
	chash1 := New(2, nil)
	chash2 := New(2, nil)
	a, b, c := &serviceNode{"testnode_aa"}, &serviceNode{"testnode_bb"}, &serviceNode{"testnode_cc"}
	chash1.AddNodes(a, b, c)
	chash2.AddNodes(a, b, c)
	for i := 0; i < 1000; i++ {
		key := myrand.Uint32()
		if chash1.GetNodeByKey(key) != chash2.GetNodeByKey(key) {
			t.Error("fail key:", key, chash1.GetNodeByKey(key).Name(), chash2.GetNodeByKey(key).Name())
		}
	}

	// chash1.AddNodes(&serviceNode{"testnode_aa"}, &serviceNode{"testnode_bb"}, &serviceNode{"testnode_cc"})
	// chash2.AddNodes(&serviceNode{"testnode_aa"}, &serviceNode{"testnode_bb"}, &serviceNode{"testnode_cc"})
	// for i := 0; i < 1000; i++ {
	// 	key := myrand.Uint32()
	// 	if chash1.GetNodeByKey(key).Name() != chash2.GetNodeByKey(key).Name() {
	// 		t.Error("fail key:", key, chash1.GetNodeByKey(key).Name(), chash2.GetNodeByKey(key).Name())
	// 	}
	// }
}
func BenchmarkGetNodeByKey(b *testing.B) {
	ch := New(4, nil)
	nodes := make([]Node, 4)
	for i, _ := range nodes {
		nodes[i] = &serviceNode{strconv.Itoa(i) + "testnode"}
	}
	ch.AddNodes(nodes...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch.GetNodeByKey(uint32(i))
	}
}
