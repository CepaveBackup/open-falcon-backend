package client

import (
	log "github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("INFO")

func SetLoggerLevel(level string) {
	logger = log.NewDefaultLogger(level)
}
