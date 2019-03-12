package web

import (
	"derek82511/jt/config"
	"derek82511/jt/service/log"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	requestLogger "github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func SetupRecover(app *iris.Application) {
	app.Use(recover.New())
}

func SetupRequestLogger(app *iris.Application) {
	app.Use(requestLogger.New(requestLogger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Query:   false,
		Columns: false,
		LogFunc: func(endTime time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
			endTimeFormatted := endTime.Format("2006/01/02 - 15:04:05")
			output := requestLogger.Columnize(endTimeFormatted, latency, status, ip, method, path, message, headerMessage)

			log.WebLogger.Info(strings.Split(output, "\n")[1])
		},
	}))
}

func SetupSite(app *iris.Application) {
	reporthandlers := &app.StaticWeb("/jmeter/reports", config.JMETER_REPORTS_FOLDER).Handlers
	*reporthandlers = append(context.Handlers{func(ctx iris.Context) {
		ctx.Request().URL.Path = strings.Replace(ctx.Request().URL.Path, "/index.html", "/main.html", -1)
		ctx.Next()
	}}, (*reporthandlers)...)

	spa := app.StaticHandler(config.JMETER_SITE_FOLDER, false, false)
	app.SPA(spa)
}
