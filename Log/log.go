package Log

import (
	"github.com/natefinch/lumberjack"
	"github.com/uniplaces/carbon"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var LogStash *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		// 日志文件的位置
		Filename: "./logs/logstash-" + carbon.Now().DateString() + ".log",
		// 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxSize: 10,
		// 保留旧文件的最大个数
		//MaxBackups: 5,
		// 保留旧文件的最大天数
		MaxAge: 30,
		// 是否压缩/归档旧文件
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	LogStash = logger.Sugar()

	defer func() {
		_ = LogStash.Sync()
	}()

}

func simpleHttpGet(url string) {
	LogStash.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		LogStash.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		LogStash.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		defer func() {
			_ = resp.Body.Close()
		}()
	}
}
