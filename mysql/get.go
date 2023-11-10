package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
)

const BatchSize = 100 * 100
const BatchSize5W = 5 * BatchSize

func BatchGet(client *gorm.DB, tableName string, conditions []Condition, keys ...string) (res []map[string]any, err error) {
	var count int64
	var wg sync.WaitGroup
	var lock sync.Mutex

	selSql := `select `
	countSql := `select count(*) `
	if len(keys) != 0 {
		for i, key := range keys {
			selSql += " " + key + " "
			if i != len(keys)-1 {
				selSql += ", "
			}
		}
	} else {
		selSql += " * "
	}
	condition := " from " + tableName + " where 1 = 1 "
	for _, cond := range conditions {
		condition += AddConditions(cond.Symbol, cond.ColumnName, cond.ColumnValue)
	}
	selSql += condition
	countSql += condition
	err = client.Raw(countSql).Count(&count).Error
	if err != nil {
		return
	}

	res = make([]map[string]any, 0, count)
	var times = int(count) / BatchSize
	if int(count)%BatchSize > 0 {
		times++
	}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(i int) {
			var temp = make([]map[string]any, 0, BatchSize)
			start, end := (i)*BatchSize, (i+1)*BatchSize
			sql := selSql + fmt.Sprintf(` limit %d,%d`, start, end)
			err2 := client.Raw(sql).Scan(&temp).Error
			if err2 != nil {
				panic(err)
			}
			lock.Lock()
			res = append(res, temp...)
			lock.Unlock()
			wg.Done()
			log.Printf("times %d is done!", i)

		}(i)
	}
	wg.Wait()
	return

}
