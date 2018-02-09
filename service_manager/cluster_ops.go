package servicemanager

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/eventing/logging"
	"github.com/couchbase/eventing/util"
)

var getEventingNodesAddressesOpCallback = func(args ...interface{}) error {
	m := args[0].(*ServiceMgr)

	hostAddress := net.JoinHostPort(util.Localhost(), m.restPort)

	eventingNodeAddrs, err := util.EventingNodesAddresses(m.auth, hostAddress)
	if err != nil {
		logging.Errorf("SMCO Failed to get all eventing nodes, err: %v", err)
		return err
	} else if len(eventingNodeAddrs) == 0 {
		logging.Errorf("SMCO Count of eventing nodes reported is 0, unexpected")
		return fmt.Errorf("eventing node count reported as 0")
	} else {
		logging.Debugf("SMCO Got eventing nodes: %#v", eventingNodeAddrs)
		m.eventingNodeAddrs = eventingNodeAddrs
		return nil
	}

}

var getHTTPServiceAuth = func(args ...interface{}) error {
	m := args[0].(*ServiceMgr)

	clusterURL := net.JoinHostPort(util.Localhost(), m.restPort)
	user, password, err := cbauth.GetHTTPServiceAuth(clusterURL)
	if err != nil {
		logging.Errorf("SMCO Failed to get cluster auth details, err: %v", err)
		return err
	}

	m.auth = fmt.Sprintf("%s:%s", user, password)
	return nil
}

var storeKeepNodesCallback = func(args ...interface{}) error {
	keepNodeUUIDs := args[0].([]string)

	data, err := json.Marshal(&keepNodeUUIDs)
	if err != nil {
		logging.Errorf("SMCO Failed to marshal keepNodes: %v, err: %v",
			keepNodeUUIDs, err)
		return err
	}

	err = util.MetakvSet(metakvConfigKeepNodes, data, nil)
	if err != nil {
		logging.Errorf("SMCO Failed to store keep nodes UUIDs: %v in metakv, err: %v",
			keepNodeUUIDs, err)
		return err
	}

	return nil
}

var stopRebalanceCallback = func(args ...interface{}) error {
	r := args[0].(*rebalancer)

	logging.Errorf("SMCO Updating metakv to signify rebalance cancellation")

	path := metakvRebalanceTokenPath + r.change.ID
	err := util.MetakvSet(path, []byte(stopRebalance), nil)
	if err != nil {
		logging.Errorf("SMCO Failed to update rebalance token: %v in metakv as part of cancelling rebalance, err: %v",
			r.change.ID, err)
		return err
	}

	return nil
}

var cleanupEventingMetaKvPath = func(args ...interface{}) error {
	path := args[0].(string)

	err := util.RecursiveDelete(path)
	if err != nil {
		logging.Errorf("Failed to purge eventing artifacts from path: %v, err: %v", path, err)
	}

	return err
}

var metaKVSetCallback = func(args ...interface{}) error {
	path := args[0].(string)
	changeID := args[1].(string)

	err := util.MetakvSet(path, []byte(changeID), nil)
	if err != nil {
		logging.Errorf("Failed to store into metakv path: %v, err: %v", path, err)
	}

	return err
}
