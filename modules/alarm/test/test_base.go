package test

import (
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/alarm/g"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/alarm/model"
)

func initTest() {
	g.ParseConfig("../test_cfg.json")
	model.InitDatabase()
}
