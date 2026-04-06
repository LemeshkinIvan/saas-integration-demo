package logger

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logg      *zap.Logger
	asyncChan = make(chan LogMessage, 1000)
)

type LogMessage struct {
	Level   zapcore.Level
	Message string
	Fields  []zap.Field
}

// InitLogger и запуск асинхронного логгера
func InitLogger() {
	// настройка ротации логов
	rotator := &lumberjack.Logger{
		Filename:   "../logs/app.log",
		MaxSize:    5, // MB
		MaxBackups: 12,
		MaxAge:     30, // days
		Compress:   true,
	}

	// кодеки (JSON или console)
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)

	// комбинированный writer: stdout + файл
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(rotator), zapcore.InfoLevel),
	)

	Logg = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// прокидываем gin логгер
	gin.DefaultWriter = zapcore.AddSync(os.Stdout)
	gin.DefaultErrorWriter = zapcore.AddSync(os.Stdout)

	startAsyncLogger()
}

// Асинхронное логгирование
func startAsyncLogger() {
	go func() {
		for msg := range asyncChan {
			switch msg.Level {
			case zapcore.DebugLevel:
				Logg.Debug(msg.Message, msg.Fields...)
			case zapcore.InfoLevel:
				Logg.Info(msg.Message, msg.Fields...)
			case zapcore.WarnLevel:
				Logg.Warn(msg.Message, msg.Fields...)
			case zapcore.ErrorLevel:
				Logg.Error(msg.Message, msg.Fields...)
			case zapcore.FatalLevel:
				Logg.Fatal(msg.Message, msg.Fields...)
			}
		}
	}()
}

// Асинхронный вызов логирования
func AsyncLog(level zapcore.Level, msg string, fields ...zap.Field) {
	select {
	case asyncChan <- LogMessage{Level: level, Message: msg, Fields: fields}:
	default:
		// fallback — если очередь переполнена, пишем синхронно
		Logg.WithOptions(zap.AddCallerSkip(1)).Log(level, msg, fields...)
	}
}

// Остановка логгера (graceful shutdown)
func StopAsyncLogger(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		close(asyncChan)
		Logg.Sync()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		Logg.Warn("Async logger shutdown timeout")
	}
}
