package main

import (
	"os"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
	. "github.com/sniperkit/snk.golang.vuejs-multi-backend/logger"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"

	// external
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() {
	db, err := gorm.Open(config.Global.Store.Dialect, config.Global.Store.DSN)
	if err != nil {
		Error("open db connect error.")
		os.Exit(-1)
	}

	db.DB().SetMaxIdleConns(config.Global.Store.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.Global.Store.MaxOpenConns)

	// global
	model.DB = db
}
