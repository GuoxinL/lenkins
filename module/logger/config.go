package logger

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLog 初始化日志 logger
func InitLog(logPath string, logLevel zapcore.Level) {
	config := zapcore.EncoderConfig{
		MessageKey:   "msg",                            //结构化（json）输出：msg的key
		LevelKey:     "level",                          //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                             //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                           //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalColorLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,       //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		//输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}
	// 获取io.Writer的实现
	infoWriter := GetWriter(logPath)
	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), logLevel),                             //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel)).Named("Lenkins")
	zap.ReplaceGlobals(logger)
	err := logger.Sync()
	if err != nil {
		return
	}

}

func GetWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,    //最大M数，超过则切割
		MaxBackups: 5,     //最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30,    //保存30天
		Compress:   false, //是否压缩
	}
}
