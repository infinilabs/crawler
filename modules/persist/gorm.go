package persist

import (
	"fmt"
	log "github.com/cihub/seelog"
	api "github.com/infinitbyte/gopa/core/persist"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

type SQLORM struct {
	conn    *gorm.DB
	useLock bool
}

var dbLock sync.RWMutex

func (handler SQLORM) Get(o interface{}) error {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.First(o).Error
}

func (handler SQLORM) GetBy(field string, value interface{}, t interface{}, to interface{}) (error, api.Result) {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.Where(field+" = ?", value).First(to).Error, api.Result{}
}

func (handler SQLORM) Save(o interface{}) error {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.Save(o).Error
}

func (handler SQLORM) Create(o interface{}) error {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.Create(o).Error
}

func (handler SQLORM) Update(o interface{}) error {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.Save(o).Error
}

func (handler SQLORM) Delete(o interface{}) error {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	return handler.conn.Delete(o).Error
}

func (handler SQLORM) Count(o interface{}) (int, error) {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	var count int
	return count, handler.conn.Model(&o).Count(count).Error
}

func (handler SQLORM) GroupBy(o interface{}, selectField, groupField string, haveQuery string, haveValue interface{}) (error, map[string]interface{}) {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}
	result := map[string]interface{}{}

	start := time.Now()
	db1 := handler.conn.Model(o).Select(fmt.Sprintf("%s,count(*)", selectField)).Group(groupField)

	if haveQuery != "" {
		db1 = db1.Having(haveQuery, haveValue)
	}
	rows, err := db1.Rows()

	if err != nil {
		return err, result
	}
	type row struct {
		Field string
		Count int
	}

	for rows.Next() {
		var r row
		err = rows.Scan(&r.Field, &r.Count)
		if err != nil {
			return err, result
		}
		result[r.Field] = r.Count
	}
	log.Tracef("groupby,select:%s, group: %s, group have: %s - %v, elapsed:%vs", selectField, groupField, haveQuery, haveValue, time.Now().Sub(start).Seconds())
	return err, result
}

func (handler SQLORM) Search(t interface{}, o interface{}, q *api.Query) (error, api.Result) {
	if handler.useLock {
		dbLock.Lock()
		defer dbLock.Unlock()
	}

	if q.From < 0 {
		q.From = 0
	}
	if q.Size < 0 {
		q.Size = 10
	}

	var c int
	var err error
	start := time.Now()
	db1 := handler.conn.Model(o)
	if q.Sort != nil && len(*q.Sort) > 0 {
		for _, i := range *q.Sort {
			db1 = db1.Order(fmt.Sprintf("%s %s", i.Field, i.SortType))
		}
	}

	if q.Conds != nil {
		where := ""
		args := []interface{}{}
		for _, c1 := range q.Conds {
			if where != "" {
				where = where + " AND "
			}
			where = where + c1.Field + c1.SQLOperator + " ?"
			args = append(args, c1.Value)
		}
		log.Tracef("search where: %s - %v", where, args)
		err = db1.Limit(q.Size).Offset(q.From).Where(where, args...).Find(o).Error
		db1.Where(where, args...).Count(&c)
	} else {
		err = db1.Limit(q.Size).Offset(q.From).Find(o).Error
		db1.Count(&c)
	}

	result := api.Result{}
	result.Result = o
	result.Total = c

	log.Tracef("search,t:%v, q: %v, elapsed:%vs", t, q, time.Now().Sub(start).Seconds())
	return err, result
}
