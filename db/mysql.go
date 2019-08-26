package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ssp4599815/redis-manager/models"
)

var SqlDB *gorm.DB

func init() {
	// open a db connection
	var err error
	SqlDB, err = gorm.Open("mysql", "root:ssp123123@tcp(127.0.0.1:3306)/redis?parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	// migrate the schema
	SqlDB.AutoMigrate(&models.Node{})

}
