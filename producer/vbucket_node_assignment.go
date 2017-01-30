package producer

import (
	"fmt"
	"time"

	"github.com/couchbase/eventing/util"
	"github.com/couchbase/indexing/secondary/common"
	"github.com/couchbase/indexing/secondary/logging"
)

// Generates the vbucket to eventing node assignment, ideally generated map should
// be consistent across all nodes
func (p *Producer) vbEventingNodeAssign() {

	util.Retry(util.NewFixedBackoff(time.Second), getKVNodesAddressesOpCallback, p)

	util.Retry(util.NewFixedBackoff(time.Second), getEventingNodesAddressesOpCallback, p)

	util.Retry(util.NewFixedBackoff(time.Second), getNsServerNodesAddressesOpCallback, p)

	eventingNodeAddrs := p.getEventingNodeAddrs()
	vbucketPerNode := NumVbuckets / len(eventingNodeAddrs)
	var startVb uint16

	p.Lock()
	p.vbEventingNodeAssignMap = make(map[uint16]string)

	for i := 0; i < len(eventingNodeAddrs); i++ {
		for j := 0; j < vbucketPerNode && startVb < NumVbuckets; j++ {
			p.vbEventingNodeAssignMap[startVb] = eventingNodeAddrs[i]
			startVb++
		}
		fmt.Printf("eventing node index: %d startVb: %d\n", i, startVb)
	}
	p.Unlock()
}

func (p *Producer) getKvVbMap() {

	var cinfo *common.ClusterInfoCache

	util.Retry(util.NewFixedBackoff(time.Second), getClusterInfoCacheOpCallback, p, &cinfo)

	kvAddrs := cinfo.GetNodesByServiceType(DataService)

	p.kvVbMap = make(map[uint16]string)

	for _, kvaddr := range kvAddrs {
		addr, err := cinfo.GetServiceAddress(kvaddr, DataService)
		if err != nil {
			logging.Errorf("VBNA[%s:%d] Failed to get address of KV host, err: %v", p.AppName, p.LenRunningConsumers(), err)
			continue
		}

		vbs, err := cinfo.GetVBuckets(kvaddr, "default")
		if err != nil {
			logging.Errorf("VBNA[%s:%d] Failed to get vbuckets for given kv common.NodeId, err: %v", p.AppName, p.LenRunningConsumers(), err)
			continue
		}

		for i := 0; i < len(vbs); i++ {
			p.kvVbMap[uint16(vbs[i])] = addr
		}
	}
}
