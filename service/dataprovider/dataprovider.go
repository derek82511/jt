package dataprovider

import (
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/log"

	"github.com/jinzhu/gorm"
)

var dbInstance *gorm.DB

func GetInstance() *gorm.DB {
	if dbInstance == nil {
		dbInstance = initDb()
	}

	return dbInstance
}

func initDb() *gorm.DB {
	log.SqlLogger.Info("Initializing db ...")

	db, err := gorm.Open("sqlite3", config.JMETER_DATA_FOLDER+"/lite.db")

	if err != nil {
		log.SqlLogger.Error("failed to connect database")
		log.SqlLogger.Error(err.Error())
		panic("failed to connect database")
	}

	log.SqlLogger.Info("Success initializing db")

	db.LogMode(true)
	db.SetLogger(log.SqlLogger)

	db.AutoMigrate(&model.Scenario{}, &model.Job{})

	return db
}
