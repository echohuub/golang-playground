package main

import (
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupSimpleLogger() {
	logFileLoc, _ := os.OpenFile("simple_test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(logFileLoc)
}

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger
var customLogger *zap.Logger

func initLogger() {
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()

	encoder := getEncoder()
	writeSyncer := getLogWriter()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	customLogger = zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("custom.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	return zapcore.AddSync(file)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching url..", zap.String("url", url), zap.Error(err))
		sugarLogger.Error("[sugar]Error fetching url..", zap.String("url", url), zap.Error(err))
		customLogger.Error("[cusom]Error fetching url..", zap.String("url", url), zap.Error(err))
	} else {
		logger.Info("Success..", zap.String("statusCode", resp.Status), zap.String("url", url))
		sugarLogger.Info("[sugar]Success..", zap.String("statusCode", resp.Status), zap.String("url", url))
		customLogger.Info("[custom]Success..", zap.String("statusCode", resp.Status), zap.String("url", url))
		resp.Body.Close()
	}
}

func main() {
	// setupSimpleLogger()
	// log.Printf("hello world")
	initLogger()
	defer logger.Sync()
	simpleHttpGet("http://www.baidu.com")
	simpleHttpGet("http://www.google.com")
}
