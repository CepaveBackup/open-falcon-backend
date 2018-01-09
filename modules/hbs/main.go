package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
	oos "github.com/DistributedMonitoringSystem/open-falcon-backend/common/os"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/common/vipercfg"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/hbs/g"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/hbs/http"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/hbs/rpc"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/hbs/service"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	service.InitPackage(vipercfg.Config())
	rpc.InitPackage(vipercfg.Config())

	go http.Start()
	go rpc.Start()

	oos.HoldingAndWaitSignal(
		func(signal os.Signal) {
			rpc.Stop()
		},
		os.Interrupt, os.Kill,
		syscall.SIGTERM,
	)
}
