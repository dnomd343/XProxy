package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)
import "go.uber.org/zap/zapcore"

var Logger *zap.SugaredLogger

const (
	logTmFmt = "2006-01-02 15:04:05"
)

func init() {

	//coreConfig := zapcore.EncoderConfig{
	//	TimeKey:        "ts",
	//	LevelKey:       "level",
	//	NameKey:        "logger",
	//	CallerKey:      "caller",
	//	FunctionKey:    zapcore.OmitKey,
	//	MessageKey:     "msg",
	//	StacktraceKey:  "stacktrace",
	//	LineEnding:     zapcore.DefaultLineEnding, // 默认换行符"\n"
	//	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	//	EncodeTime:     zapcore.RFC3339TimeEncoder,    // 日志时间格式显示
	//	EncodeDuration: zapcore.MillisDurationEncoder, // 时间序列化，Duration为经过的浮点秒数
	//	EncodeCaller:   zapcore.ShortCallerEncoder,    // 日志行号显示
	//}
	//
	//encoder := zapcore.NewConsoleEncoder(coreConfig)
	//
	//newCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	//
	//logger := zap.New(newCore)
	//
	//Logger = logger.Sugar()

	GetLogger()

	//log, err := zap.NewDevelopment()
	//if err != nil {
	//	panic("log utils init failed")
	//}

	// TODO: more zap logger configure
	// TODO: reserve raw logger handle

	//Logger = log.Sugar()
}

func GetLogger() {
	config := zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    "F",
		StacktraceKey:  "S",
		SkipLineEnding: false,
		LineEnding:     "\n",

		EncodeLevel: logEncodeLevel,
		//EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     logEncodeTime,
		EncodeDuration: zapcore.StringDurationEncoder,
		//EncodeCaller:   logEncodeCaller,
		EncodeCaller: zapcore.ShortCallerEncoder,

		//EncodeName:       logEncodeName,
		ConsoleSeparator: " ",
	}

	//var zc = zap.Config{
	//	Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
	//	Development:       false,
	//	DisableCaller:     false,
	//	DisableStacktrace: false,
	//	Sampling:          nil,
	//	Encoding:          "json",
	//	EncoderConfig:     config,
	//	OutputPaths:       []string{"stdout"},
	//	ErrorOutputPaths:  []string{"stderr"},
	//	InitialFields:     map[string]interface{}{"app": "zapdex"},
	//}

	newCore := zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.Lock(os.Stderr), zapcore.DebugLevel)
	//newCore := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.Lock(os.Stderr), zapcore.DebugLevel)
	logger := zap.New(newCore, zap.AddCaller())
	//logger := zap.New(newCore)
	//logger, _ := zc.Build()
	//logger.Named("123")
	//fmt.Println(logger.Name())

	zap.ReplaceGlobals(logger)
	Logger = logger.Sugar()
}

func logEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(func(level zapcore.Level) string {
		levelStr := level.CapitalString()
		// TODO: using shell codes map
		switch level {
		case zapcore.DebugLevel:
			return fmt.Sprintf("\x1b[39m[%s]\x1b[0m", levelStr)
		case zapcore.InfoLevel:
			return fmt.Sprintf("\x1b[32m[%s]\x1b[0m", levelStr)
		case zapcore.WarnLevel:
			return fmt.Sprintf("\x1b[33m[%s]\x1b[0m", levelStr)
		case zapcore.ErrorLevel:
			return fmt.Sprintf("\x1b[31m[%s]\x1b[0m", levelStr)
		case zapcore.PanicLevel:
			return fmt.Sprintf("\x1b[95m[%s]\x1b[0m", levelStr)
		default:
			return fmt.Sprintf("[%s]", levelStr)
		}
	}(level))
}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: using `2006-01-02 15:04:05.xxx` format
	enc.AppendString(fmt.Sprintf(
		"\x1b[36m%s\x1b[0m \x1b[90m%s\x1b[0m",
		"[XProxy]", t.Format("2006-01-02 15:04:05"),
	))
}

func logEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("{" + caller.TrimmedPath() + "}")
}
