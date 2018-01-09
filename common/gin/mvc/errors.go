package mvc

import (
	ogin "github.com/DistributedMonitoringSystem/open-falcon-backend/common/gin"
	"github.com/gin-gonic/gin"
)

var NotFoundOutputBody = OutputBodyFunc(func(c *gin.Context) {
	ogin.JsonNoRouteHandler(c)
})
