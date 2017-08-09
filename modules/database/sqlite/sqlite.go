package sqlite

import (
	"fmt"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/jinzhu/gorm"
	"os"
	"path"
)

// SQLiteConfig currently do nothing
type SQLiteConfig struct {
}

// GetInstance return sqlite instance for further access
func GetInstance(cfg *SQLiteConfig) *gorm.DB {
	os.MkdirAll(path.Join(global.Env().SystemConfig.GetWorkingDir(), "database/"), 0777)
	fileName := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_busy_timeout=50000000", path.Join(global.Env().SystemConfig.GetWorkingDir(), "database/db.sqlite"))

	var err error
	db, err := gorm.Open("sqlite3", fileName)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
