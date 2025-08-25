package web

import (
	"derek82511/jt/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func SetupRecover(app *iris.Application) {
	app.Use(recover.New())
}

func SetupRequestLogger(app *iris.Application) {
	app.Use(logger.New())
}

func SetupSite(app *iris.Application) {
	app.HandleDir("/jmeter/reports", config.JMETER_REPORTS_FOLDER)

	// Serve the Angular frontend's static files with SPA support.
	// This will serve "index.html" when a file is not found.
	app.HandleDir("/", config.JMETER_SITE_FOLDER, iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
}
