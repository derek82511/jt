package log

import (
	"bufio"
	"derek82511/jt/config"
	"os"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var WebLogger *logrus.Logger
var SqlLogger *logrus.Logger

func init() {
	Logger = initLogger("ap")
	WebLogger = initLogger("web")
	SqlLogger = initLogger("sql")
}

func initLogger(name string) *logrus.Logger {
	path := config.JMETER_LOG_FOLDER + "/" + name + ".log"

	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(time.Duration(1)*time.Hour),
	)

	if err != nil {
		logrus.Errorf("logger error. %v", errors.WithStack(err))
	}

	newLogger := logrus.New()

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.Errorf("logger error when set null. %v", errors.WithStack(err))
	}

	newLogger.SetOutput(bufio.NewWriter(src))

	newLogger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.TextFormatter{},
	))

	return newLogger
}
