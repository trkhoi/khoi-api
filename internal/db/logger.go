package db

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type LogrusLogger struct {
	logger   *logrus.Entry
	logLevel logger.LogLevel
}

func NewLogrusLogger(l *logrus.Entry) *LogrusLogger {
	gormLockLvl := logger.Error
	switch l.Logger.Level {
	case logrus.DebugLevel:
		gormLockLvl = logger.Info
	case logrus.WarnLevel:
		gormLockLvl = logger.Warn
	}

	return &LogrusLogger{
		logger:   l,
		logLevel: gormLockLvl,
	}
}

// LogMode sets the logging mode dynamically
func (l *LogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level // set the logging level
	return l
}

// Info print infomation and/or useful tips.
func (l *LogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		fields := logrus.Fields{"context": ctx, "data": data}
		l.logger.WithFields(fields).Info(msg)
	}
}

// Warn print warning messages.
func (l *LogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		fields := logrus.Fields{"context": ctx, "data": data}
		l.logger.WithFields(fields).Warn(msg)
	}
}

// Error print error messages and stack trace
func (l *LogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		fields := logrus.Fields{"context": ctx, "data": data}
		l.logger.WithFields(fields).Error(msg)
	}
}

// Trace print sql message
func (l *LogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= logger.Error:
		sql, rows := fc()
		l.logger.WithFields(logrus.Fields{
			"exec_at":       begin.Format("2006-01-02 15:04:05"),
			"elapsed_time":  elapsed,
			"sql":           sql,
			"rows_affected": rows,
		}).WithError(err).Error("GORM Error")
	case elapsed > time.Duration(500)*time.Millisecond && l.logLevel >= logger.Warn:
		sql, rows := fc()
		l.logger.WithFields(logrus.Fields{
			"exec_at":       begin.Format("2006-01-02 15:04:05"),
			"elapsed_time":  elapsed,
			"sql":           sql,
			"rows_affected": rows,
		}).Warn("GORM Slow SQL")
	case l.logLevel >= logger.Info:
		sql, rows := fc()
		l.logger.WithFields(logrus.Fields{
			"exec_at":       begin.Format("2006-01-02 15:04:05"),
			"elapsed_time":  elapsed,
			"sql":           sql,
			"rows_affected": rows,
		}).Info("GORM SQL")
	}
}
