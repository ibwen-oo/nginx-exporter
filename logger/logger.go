package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Log *zap.Logger
)

// zapcore.NewTee 指定日志输出到多个地方
func Init(logPath string) (err error) {

	/* 定义 core 变量 */
	var cores zapcore.Core
	/* 定义日志输出到文件的encoder */
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	writeSync := zapcore.AddSync(file)
	// 修改时间格式,添加日志级别(大写)
	// 构造配置信息
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encodeConf)

	/* 定义日志输出到终端的encoder */
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	/* 配置日志输出到多端 */
	cores = zapcore.NewTee(
		// 日志输出到文件
		zapcore.NewCore(encoder, writeSync, zapcore.ErrorLevel),
		// 日志输出到终端
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
	)

	lg := zap.New(cores, zap.AddCaller())
	// 替换zap包中全局的logger实例,后续在其它包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(lg)
	return
}