package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

func NewLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetReportCaller(true)

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "Jan 02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				repopath := "C:/Users/Christopher Robin/Documents/Coding/Golang/Nusastra/" // JANGAN LUPA GANTI PAAKE ENV
				filename := strings.Replace(f.File, repopath, " ", 1)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			ForceColors:      true,
			QuoteEmptyFields: true,
		})
		// JANGAN LUPA GANTI PATH PAAKE ENV
		logFileName := filepath.Join("C:/Users/Christopher Robin/Documents/Coding/Golang/Nusastra/storage/logs", fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open log file: %v", err))
		}

		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(multiWriter)

		if os.Getenv("APP_ENV") == "production" {
			logger.SetLevel(logrus.InfoLevel)
		} else {
			logger.SetLevel(logrus.InfoLevel)
		}
	})

	return logger
}

func ErrorWithTraceID(fields map[string]interface{}, msg string) uuid.UUID {
	traceID, err := uuid.NewRandom()
	if err != nil {
		Error(map[string]interface{}{
			"error": err.Error(),
		}, "[log.ErrorWithTraceID] failed to generate trace ID")
	}

	fields["trace_id"] = traceID
	logger.WithFields(fields).Error(msg)

	return traceID
}

func Debug(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Debug(msg)
}

func Info(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Info(msg)
}

func Warn(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Warn(msg)
}

func Error(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Error(msg)
}

func Fatal(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Fatal(msg)
}

func Panic(fields map[string]interface{}, msg string) {
	logger.WithFields(fields).Panic(msg)
}
