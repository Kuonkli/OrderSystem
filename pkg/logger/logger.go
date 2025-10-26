package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Entry
}

func Init(serviceName string) *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	return &Logger{
		Entry: logger.WithField("service", serviceName),
	}
}

func (l *Logger) WithRequestID(id string) *logrus.Entry {
	return l.Entry.WithField("request_id", id)
}
