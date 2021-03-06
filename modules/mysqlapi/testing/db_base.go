package testing

import (
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb"
	check "gopkg.in/check.v1"
)

// The base environment for RDB testing
func InitRdb(c *check.C) {
	dbTest.SetupByViableDbConfig(c, rdb.InitPortalRdb)
}
func ReleaseRdb(c *check.C) {
	rdb.ReleaseAllRdb()
}
