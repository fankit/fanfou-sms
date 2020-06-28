package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
	"tsmsrv/conf"
)

type LoggerMgr struct {
	Log *zap.Logger
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

var Logger *LoggerMgr //*zap.Logger

func NewLoggerMgr(levels string) (err error) {
	var (
		maxSize    int
		maxBuckups int
		maxAge     int
	)

	if maxSize, err = conf.GlobConfig.LogSection().Key(`max_size`).Int(); err != nil {
		log.Fatal(`max_siz`, err.Error())
		return
	}

	if maxBuckups, err = conf.GlobConfig.LogSection().Key(`max_backups`).Int(); err != nil {
		log.Fatal(`max_backup`, err.Error())
		return
	}

	if maxAge, err = conf.GlobConfig.LogSection().Key(`max_age`).Int(); err != nil {
		log.Fatal(`max_age`, err.Error())
		return
	}

	hook := lumberjack.Logger{
		Filename:   conf.GlobConfig.LogSection().Key(`acc_path`).String(),
		MaxSize:    maxSize,
		MaxBackups: maxBuckups,
		MaxAge:     maxAge,
		Compress:   false,
	}

	zws := zapcore.AddSync(&hook)

	var level zapcore.Level

	switch levels {
	case "info":
		level = zap.InfoLevel
	case "debug":
		level = zap.DebugLevel
	case "err":
		level = zap.ErrorLevel
	case "fata":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zws, level))
	Logger = &LoggerMgr{Log: logger}
	return
}

func (l *LoggerMgr) Sync(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			l.Log.Sync()
		}
	}()
}

func (l *LoggerMgr) Color(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func (l *LoggerMgr) HttpMethodColor(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
