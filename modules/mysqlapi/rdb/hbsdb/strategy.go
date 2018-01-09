package hbsdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/DistributedMonitoringSystem/open-falcon-backend/common/model"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/set"
)

// 获取所有的Strategy列表
func QueryStrategies(tpls map[int]*model.NewTemplate) (map[int]*model.NewStrategy, error) {
	ret := make(map[int]*model.NewStrategy)

	if tpls == nil || len(tpls) == 0 {
		return ret, fmt.Errorf("illegal argument")
	}

	now := time.Now().Format("15:04")
	sql := fmt.Sprintf(
		"select %s from strategy as s where (s.run_begin='' and s.run_end='') or (s.run_begin <= '%s' and s.run_end > '%s')",
		"s.id, s.metric, s.tags, s.func, s.op, s.right_value, s.max_step, s.priority, s.note, s.tpl_id",
		now,
		now,
	)

	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		s := model.NewStrategy{}
		var tags string
		var tid int
		err = rows.Scan(&s.ID, &s.Metric, &tags, &s.Func, &s.Operator, &s.RightValue, &s.MaxStep, &s.Priority, &s.Note, &tid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		tt := make(map[string]string)

		if tags != "" {
			arr := strings.Split(tags, ",")
			for _, tag := range arr {
				kv := strings.Split(tag, "=")
				if len(kv) != 2 {
					continue
				}
				if strings.Contains(kv[1], "$") {
					strategyId := s.ID
					tagName := kv[1]
					sql = fmt.Sprintf(
						"select %s from tags where strategy_id = %d and name = '%s'",
						"value",
						strategyId,
						tagName,
					)
					row := DB.QueryRow(sql)
					if err != nil {
						log.Println("ERROR:", err.Error())
						continue
					}
					var value string
					err = row.Scan(&value)
					if err != nil {
						log.Println("ERROR:", err.Error())
						continue
					}
					// log.Println("templateId =", tid)
					// log.Println("strategyId =", strategyId)
					// log.Println(kv[1] + " =", value)
					kv[1] = value
				}
				tt[kv[0]] = kv[1]
			}
		}

		s.Tags = tt
		s.Tpl = tpls[tid]
		if s.Tpl == nil {
			log.Printf("WARN: tpl is nil. strategy id=%d, tpl id=%d", s.ID, tid)
			// 如果Strategy没有对应的Tpl，那就没有action，就没法报警，无需往后传递了
			continue
		}

		ret[s.ID] = &s
	}

	return ret, nil
}

func QueryBuiltinMetrics(tids string) ([]*model.NewBuiltinMetric, error) {
	sql := fmt.Sprintf(
		"select metric, tags from strategy where tpl_id in (%s) and metric in ('net.port.listen', 'proc.num', 'du.bs', 'url.check.health')",
		tids,
	)

	ret := []*model.NewBuiltinMetric{}

	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	metricTagsSet := set.NewStringSet()

	defer rows.Close()
	for rows.Next() {
		builtinMetric := model.NewBuiltinMetric{}
		err = rows.Scan(&builtinMetric.Metric, &builtinMetric.Tags)
		if err != nil {
			log.Println("WARN:", err)
			continue
		}

		k := fmt.Sprintf("%s%s", builtinMetric.Metric, builtinMetric.Tags)
		if metricTagsSet.Exists(k) {
			continue
		}

		ret = append(ret, &builtinMetric)
		metricTagsSet.Add(k)
	}

	return ret, nil
}
