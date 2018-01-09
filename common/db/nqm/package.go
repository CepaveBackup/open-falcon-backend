package nqm

import (
	f "github.com/DistributedMonitoringSystem/open-falcon-backend/common/db/facade"
	log "github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
	tb "github.com/DistributedMonitoringSystem/open-falcon-backend/common/textbuilder"
)

var DbFacade *f.DbFacade

var t = tb.Dsl

var logger = log.NewDefaultLogger("warn")
