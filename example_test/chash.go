package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	chash "github.com/jursonmo/consistenthash"
)

type serviceNode struct {
	name string
}

func (n *serviceNode) Name() string {
	return n.name
}

func main() {
	myrand := rand.New(rand.NewSource(time.Now().Unix()))
	ch := chash.New(2, nil)
	aanode := &serviceNode{name: "testnode_aa"}
	ch.AddNodes(aanode)
	ch.ShowRing()
	fmt.Println("==========================")
	bbnode := &serviceNode{name: "testnode_bb"}
	ccnode := &serviceNode{name: "testnode_cc"}
	ch.AddNodes(bbnode, ccnode)
	ch.ShowRing()
	fmt.Println("==========================")
	ch.DelNode(aanode)
	ch.ShowRing()

	fmt.Println("==========================")
	pre_rid := uint32(0)
	rid := uint32(0)
	rlen := ch.GetRingLen()
	for i := 0; i < rlen; i++ {
		rid = ch.GetRingNodeId(i)
		rand_rid := pre_rid + uint32(myrand.Int31n(int32(rid-pre_rid))) //rand_rid is (pre_rid, rid]
		node1 := ch.GetNodeByKey(rand_rid)
		node2 := ch.GetNodeByKey(rid)
		fmt.Printf("rid=%d, rand_rid=%d, %s,%s\n", rid, rand_rid, node1.Name(), node2.Name())
		if node1 != node2 {
			fmt.Printf("warning: rid=%d, rand_rid=%d, %s,%s\n", rid, rand_rid, node1.Name(), node2.Name())
		}
		pre_rid = rid
	}
	rid = ch.GetRingNodeId(rlen-1) + 1
	rand_rid := rid + uint32(myrand.Int31n(int32(math.MaxUint32-rid)))
	node1 := ch.GetNodeByKey(rid)
	node2 := ch.GetNodeByKey(rand_rid)
	fmt.Printf(" rid=%d, rand_rid=%d, %s,%s\n", rid, rand_rid, node1.Name(), node2.Name())
	if node1 != node2 {
		fmt.Printf("warning: rid=%d, rand_rid=%d, %s,%s\n", rid, rand_rid, node1.Name(), node2.Name())
	}
}
