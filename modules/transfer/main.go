package main

import (
	"fmt"
	"os"

	"github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/common/vipercfg"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/g"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/http"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/proc"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/receiver"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/sender"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/transfer/service"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if vipercfg.Config().GetBool("vg") {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()
	if vipercfg.Config().GetBool("debug") {
		logruslog.SetLogLevelByString("debug")
	}

	service.DefaultRelayStationFactory = service.NewRelayFactoryByGlobalConfig(g.Config())

	// proc
	proc.Start()

	sender.Start()
	receiver.Start()

	// http
	http.Start()

	select {}
}
