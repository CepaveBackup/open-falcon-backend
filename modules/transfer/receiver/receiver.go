package receiver

import (
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/receiver/rpc"
)

func Start() {
	go rpc.StartRpc()
}
