package config

import (
	log "github.com/sirupsen/logrus"
	"sync"

	cmodel "github.com/DistributedMonitoringSystem/open-falcon-backend/common/model"
	"github.com/toolkits/container/nmap"

	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/nodata/config/service"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/nodata/g"
)

// nodata配置(mockcfg)的缓存, 这些数据来自配置中心
var (
	rwlock      = sync.RWMutex{}
	NdConfigMap = nmap.NewSafeMap()
)

func Start() {
	if !g.Config().Config.Enabled {
		log.Println("config.Start warning, not enabled")
		return
	}

	service.InitDB()
	StartNdConfigCron()
	log.Println("config.Start ok")
}

// Interfaces Of StrategyMap
func SetNdConfigMap(val *nmap.SafeMap) {
	rwlock.Lock()
	defer rwlock.Unlock()

	NdConfigMap = val
}

func Keys() []string {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return NdConfigMap.Keys()
}

func Size() int {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return NdConfigMap.Size()
}

func GetNdConfig(key string) (*cmodel.NodataConfig, bool) {
	rwlock.RLock()
	defer rwlock.RUnlock()

	val, found := NdConfigMap.Get(key)
	if found && val != nil {
		return val.(*cmodel.NodataConfig), true
	}
	return &cmodel.NodataConfig{}, false
}
