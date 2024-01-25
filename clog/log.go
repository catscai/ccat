package clog

import (
	"github.com/catscai/ccat/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var gAppLogger ICatLog

func InitAppLogger() {
	cfg := config.AppCfg.LogCfg
	gAppLogger = NewZapLogger(cfg.Level, cfg.AppName, cfg.LogDir,
		cfg.MaxSize, cfg.MaxAge, cfg.MaxBackups, cfg.Console, cfg.IsFuncName)
}

func AppLogger() ICatLog {
	return gAppLogger
}

func NewZapLogger(levelStr, appName, logDir string, maxSize, maxAge, maxBackups int, console, isFunc bool) ICatLog {
	if !strings.HasSuffix(logDir, "/") {
		logDir += "/"
	}
	logPath := logDir + appName + ".log"
	var level zapcore.Level
	switch levelStr {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dPanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.DebugLevel
	}
	zapLogger := NewLogger(level, logPath, maxSize, maxAge, maxBackups, console, isFunc, nil)
	return &CatZapLog{
		log: zapLogger,
	}
}

func NewLogger(level zapcore.Level, logPath string, maxSize, maxAge, maxBackups int, console, isFunc bool, defaultFields map[string]interface{}) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath,    // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     maxAge,     // 文件最多保存多少天
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		Compress:   false,      // 是否压缩
		LocalTime:  true,       // 本地时间
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	if isFunc {
		encoderConfig.FunctionKey = "func"
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	// 控制台上也输出
	if console {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	var opts []zap.Option
	// 开启开发模式，堆栈跟踪// 开启文件及行号
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2), zap.Development())

	// 设置初始化字段
	fields := make([]zap.Field, 0)
	for k, v := range defaultFields {
		fields = append(fields, zap.Any(k, v))
	}
	opts = append(opts, zap.Fields(fields...))

	// 构造日志
	logger := zap.New(core, opts...)

	return logger
}

type ICatLog interface {
	Clone() ICatLog
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	DPanic(msg string, fields ...zapcore.Field)
	Panic(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
}

type CatZapLog struct {
	log *zap.Logger
}

func (zl *CatZapLog) Clone() ICatLog {
	return &CatZapLog{
		log: zl.log,
	}
}

func (zl *CatZapLog) Debug(msg string, fields ...zapcore.Field) {
	zl.log.Debug(msg, fields...)
}

func (zl *CatZapLog) Info(msg string, fields ...zapcore.Field) {
	zl.log.Info(msg, fields...)
}

func (zl *CatZapLog) Error(msg string, fields ...zapcore.Field) {
	zl.log.Error(msg, fields...)
}

func (zl *CatZapLog) Warn(msg string, fields ...zapcore.Field) {
	zl.log.Warn(msg, fields...)
}

func (zl *CatZapLog) DPanic(msg string, fields ...zapcore.Field) {
	zl.log.DPanic(msg, fields...)
}

func (zl *CatZapLog) Panic(msg string, fields ...zapcore.Field) {
	zl.log.Panic(msg, fields...)
}

func (zl *CatZapLog) Fatal(msg string, fields ...zapcore.Field) {
	zl.log.Fatal(msg, fields...)
}
