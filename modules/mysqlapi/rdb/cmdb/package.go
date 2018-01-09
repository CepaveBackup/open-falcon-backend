package cmdb

import (
	f "github.com/DistributedMonitoringSystem/open-falcon-backend/common/db/facade"
	log "github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
)

var DbFacade *f.DbFacade
var logger = log.NewDefaultLogger("warn")
