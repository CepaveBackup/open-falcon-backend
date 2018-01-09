package http

import (
	"fmt"
	"github.com/DistributedMonitoringSystem/open-falcon-backend/modules/sender/proc"
	"net/http"
)

func configProcRoutes() {

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("sms:%v, mail:%v, qq:%v", proc.GetSmsCount(), proc.GetMailCount(), proc.GetQQCount())))
	})

}
