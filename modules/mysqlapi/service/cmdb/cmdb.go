package cmdb

import (
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/mysqlapi/model"
	bossRdb "github.com/DistributedMonitoringSystem/open-falcon-backend/modules/mysqlapi/rdb/boss"
	cmdbRdb "github.com/DistributedMonitoringSystem/open-falcon-backend/modules/mysqlapi/rdb/cmdb"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/mysqlapi/service"
	bossSrv "github.com/DistributedMonitoringSystem/open-falcon-backend/modules/mysqlapi/service/boss"
)

func SyncDataFromBoss() (*model.OwlScheduleLog, error) {
	return service.ScheduleService.Execute(
		model.NewSchedule("import.imdb", 300),
		func() error {
			sourceData := bossSrv.Boss2cmdb(bossRdb.GetSyncData())

			cmdbRdb.SyncForHosts(sourceData)

			return nil
		},
	)
}
