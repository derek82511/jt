package dataprovider

import (
	"derek82511/jt/config"
	"derek82511/jt/model"
	"derek82511/jt/service/log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.SqlLogger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open(config.JMETER_DATA_FOLDER+"/lite.db"), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.SqlLogger.Error("failed to connect database")
		log.SqlLogger.Error(err.Error())
		panic("failed to connect database")
	}

	log.SqlLogger.Info("Success initializing db")

	db.AutoMigrate(&model.Scenario{}, &model.Job{})

	return db
}
