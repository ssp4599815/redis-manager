package main

import (
	"github.com/ssp4599815/redis-manager/db"
)

func main() {
	defer db.SqlDB.Close()
	r := initRouter()
	r.Run(":8089")
}
