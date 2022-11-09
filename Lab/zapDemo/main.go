package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

//logger日志记录器
//var logger *zap.Logger
//
//func main() {
//	InitLogger()
//	defer logger.Sync()
//	simpleHttpGet("www.baidu.com")
//	simpleHttpGet("http://www.baidu.com")
//}
//
//func InitLogger() {
//	logger, _ = zap.NewProduction()
//}
//
//func simpleHttpGet(url string) {
//	resp, err := http.Get(url)
//	if err != nil {
//		logger.Error(
//			"Error fetching url..",
//			zap.String("url", url),
//			zap.Error(err))
//	} else {
//		logger.Info("Success..",
//			zap.String("statusCode", resp.Status),
//			zap.String("url", url))
//		resp.Body.Close()
//	}
//}

//Sugared Logger日志记录器********************************************************
//var sugarLogger *zap.SugaredLogger
//
//func main() {
//	InitLogger()
//	defer sugarLogger.Sync()
//	simpleHttpGet("www.baidu.com")
//	simpleHttpGet("http://www.baidu.com")
//}
//
//func InitLogger() {
//	logger, _ := zap.NewProduction()
//	sugarLogger = logger.Sugar()
//}
//
//func simpleHttpGet(url string) {
//	sugarLogger.Debugf("Trying to hit GET request for %s", url)
//	resp, err := http.Get(url)
//	if err != nil {
//		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
//	} else {
//		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
//		resp.Body.Close()
//	}
//}

//定制Logger：将日志写入文件而不是终端，时间终端
//
//var sugarLogger *zap.SugaredLogger
//
//func main() {
//	InitLogger()
//	defer sugarLogger.Sync()
//	simpleHttpGet("www.google.com")
//	simpleHttpGet("http://www.google.com")
//}
//func InitLogger() {
//	writeSyncer := getLogWriter()
//	encoder := getEncoder()
//	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
//
//	logger := zap.New(core, zap.AddCaller())
//	sugarLogger = logger.Sugar()
//}
//
//func simpleHttpGet(url string) {
//	sugarLogger.Debugf("Trying to hit GET request for %s", url)
//	resp, err := http.Get(url)
//	if err != nil {
//		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
//	} else {
//		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
//		resp.Body.Close()
//	}
//}
//
//func getEncoder() zapcore.Encoder {
//	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) //可以不按照json格式
//	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()) //不是json格式
//	//更改时间编码并添加调用者详细信息
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
//}
//
//func getLogWriter() zapcore.WriteSyncer {
//	file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//	return zapcore.AddSync(file)
//}

//使用zap接受gin框架日志********************************************
import (
	"fmt"
	"zapDemo/config"
	"zapDemo/logger"
)

func main() {
	// load config from config.json
	if len(os.Args) < 1 {
		return
	}

	if err := config.Init(os.Args[1]); err != nil {
		panic(err)
	}
	// init logger
	if err := logger.InitLogger(config.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	gin.SetMode(config.Conf.Mode)

	r := gin.Default()
	// 注册zap相关中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/hello", func(c *gin.Context) {
		// 假设你有一些数据需要记录到日志中
		var (
			name = "q1mi"
			age  = 18
		)
		// 记录日志并使用zap.Xxx(key, val)记录相关字段
		zap.L().Debug("this is hello func", zap.String("user", name), zap.Int("age", age))

		c.String(http.StatusOK, "hello liwenzhou.com")
	})

	addr := fmt.Sprintf(":%v", config.Conf.Port)
	r.Run(addr)
}
