/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlite

import (
	"fmt"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

type UserInfo struct {
	Uid        int    `gorm:"AUTO_INCREMENT"`
	Count      int    `gorm:"-"`
	Username   string `gorm:"size:255"`
	DepartName string `gorm:"size:255"`
	Created    time.Time
}

type UserGroup struct {
	Count      int
	DepartName string
}

func TestSmokeTest1(t *testing.T) {
	util.FileDelete("/tmp/test_database12.db")

	fileName := fmt.Sprintf("file:%s?cache=shared&mode=rwc", "/tmp/test_database12.db")
	fmt.Println(fileName)

	db, err := gorm.Open("sqlite3", fileName)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&UserInfo{})

	u := UserInfo{Username: "medcl", DepartName: "dev"}
	db.Create(&u)

	u = UserInfo{Username: "shay", DepartName: "dev"}
	db.Create(&u)

	u = UserInfo{Username: "joe", DepartName: "design"}
	db.Create(&u)

	rows, _ := db.Table("user_infos").Select("depart_name,count(*) as count").Group("depart_name").
		Having("username=?", "medcl").
		Rows()
	columns, _ := rows.Columns()
	fmt.Println(columns)

	g := UserGroup{}

	for rows.Next() {
		db.ScanRows(rows, &g)
		fmt.Println(g)
	}

	db.AutoMigrate(model.Host{})
	host := model.Host{}
	host.Host = "baidu.com"
	time := time.Now().UTC()
	host.Created = &time
	host.Updated = &time

	db.Create(&host)
	host = model.Host{}
	host.Host = "baidu.com"
	db.Find(&host)
	fmt.Println(util.ToJson(host, true))

	var us []UserInfo
	db.Model(&u).Where("depart_name=?", "dev").Find(&us)
	fmt.Println(us)

}
