package cron

import (
	log "github.com/DistributedMonitoringSystem/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("WARN")

func SetLoggerLevel(level string) {
	logger = log.NewDefaultLogger(level)
}
