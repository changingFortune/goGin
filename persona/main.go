package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lew/persona/app/config"
	"github.com/lew/persona/app/model"
	"github.com/lew/persona/app/model/basemgo"
	"github.com/lew/persona/app/model/redis"
	"github.com/lew/persona/app/router"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func initLog(cfgErr error) {
	// 设置日志格式为json格式　自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}

	// log.SetFormatter(&log.TextFormatter{
	// 	FullTimestamp: true,
	// })
	log.SetReportCaller(true)

	logLevel := log.WarnLevel

	if cfgErr == nil {

		commonPath := config.Cfg.MustValue("log", "dir", "./logs") + "/gin_common"
		fatalPath := config.Cfg.MustValue("log", "dir", "./logs") + "/gin_err"
		commonWriter, _ := rotatelogs.New(
			commonPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(commonPath),                            // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(24*time.Hour),                            // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second), // 日志切割时间间隔
		)
		errWriter, _ := rotatelogs.New(
			fatalPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(fatalPath),
			rotatelogs.WithMaxAge(2*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		log.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				log.DebugLevel: commonWriter, // 为不同级别设置不同的输出目的
				log.InfoLevel:  commonWriter,
				log.WarnLevel:  commonWriter,
				log.ErrorLevel: errWriter,
				log.FatalLevel: errWriter,
				log.PanicLevel: errWriter,
			},
			&log.TextFormatter{
				FullTimestamp: true,
			},
		))
		// fileInfo, _ := os.Stat(file)
		// if fileInfo == nil {
		// 	os.Mkdir("./logs", os.ModePerm)
		// 	f, _ = os.Create(file) // 文件不存在就创建
		// } else {
		// 	f, _ = os.OpenFile(file, os.O_RDWR, 0666) // 文件存在就打开 os.ModePerm
		// }

		//日志级别
		level := config.Cfg.MustValue("log", "outLevel", "DebugLevel")
		switch level {
		case "TraceLevel":
			logLevel = log.WarnLevel
		case "DebugLevel":
			logLevel = log.DebugLevel
		case "InfoLevel":
			logLevel = log.InfoLevel
		case "WarnLevel":
			logLevel = log.WarnLevel
		case "ErrorLevel":
			logLevel = log.ErrorLevel
		case "FatalLevel":
			logLevel = log.FatalLevel
		case "PanicLevel":
			logLevel = log.PanicLevel

		}

	}
	// else {
	// 	// // 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 	// // 日志消息输出可以是任意的io.writer类型
	// f := os.Stdout
	// 	log.SetOutput(f)
	// }

	// 设置sentry
	hook, err := logrus_sentry.NewSentryHook("http://ac69e21b51184f49b36b440333ed9e72:2236075a82924a599493b38747c7b056@e.datangyouxi.com/16", []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err == nil {
		hook.Timeout = 5 * time.Second
		hook.StacktraceConfiguration.Enable = true
		hook.SetEnvironment(*env)
		log.AddHook(hook)
	}

	// 设置日志级别为warn以上
	log.SetLevel(logLevel)

	if cfgErr != nil {
		log.Fatal("initLog Fatal:", cfgErr.Error())
	}
	log.Debug("debug test")
}
func initSentry() {
	raven.SetDSN("http://ac69e21b51184f49b36b440333ed9e72:2236075a82924a599493b38747c7b056@e.datangyouxi.com/16")
}

var env *string = flag.String("env", "debug", "debug alpha release")

var ports *string = flag.String("ports", "9008", "9008,9009,9010,9011")

func main() {
	// 获取环境参数
	flag.Parse()
	fmt.Println("env:%s ports:%s", *env, *ports)
	initSentry()
	cfgErr := config.InitCfg(env)
	initLog(cfgErr)
	redis.InitRedis()
	basemgo.InitMongo()
	model.DailyClean()
	fmt.Println("main main")
	router.Register(env, ports)

}
