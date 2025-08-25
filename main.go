package main

import (
	"derek82511/jt/config"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/log"
	"derek82511/jt/service/web"
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
)

func main() {
	// init db instance
	log.WebLogger.Info("Initializing db ...")
	db := dataprovider.GetInstance()
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// make data folders
	os.MkdirAll(config.JMETER_SCRIPTS_FOLDER, os.ModePerm)
	os.MkdirAll(config.JMETER_REPORTS_FOLDER, os.ModePerm)

	// init web server
	log.WebLogger.Info("Initializing web server ...")

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Logger().SetOutput(log.WebLogger.Out)

	web.SetupRecover(app)
	web.SetupRequestLogger(app)
	web.SetupSite(app)
	web.SetupApi(app)
	web.SetupWebsocket(app)

	log.WebLogger.Info("Web server is running at ", config.JMETER_SERVICE_PORT)
	err = app.Run(iris.Addr(":" + config.JMETER_SERVICE_PORT))
	if err != nil {
		// Use unbuffered write to stderr to guarantee the message is seen before exit.
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
