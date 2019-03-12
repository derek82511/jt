package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"

	"derek82511/jt/config"
	"derek82511/jt/service/dataprovider"
	"derek82511/jt/service/web"
)

func main() {
	db := dataprovider.GetInstance()
	defer db.Close()

	app := iris.New()

	web.SetupRecover(app)
	web.SetupRequestLogger(app)
	web.SetupSite(app)
	web.SetupApi(app)
	web.SetupWebsocket(app)

	app.Run(iris.Addr(":" + config.JMETER_SERVICE_PORT))
}
