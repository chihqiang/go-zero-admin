package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
)

type Logger struct {
	LogLevel logger.LogLevel
}

// LogMode 设置日志级别，返回新的 Logger 实例
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 输出信息日志
func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

// Warn 输出警告日志
func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

// Error 输出错误日志
func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(s, i...)
}

// Trace 输出 SQL 追踪日志
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 根据是否有错误使用不同的日志级别
	if err != nil && l.LogLevel >= logger.Error {
		logx.WithContext(ctx).WithDuration(elapsed).Errorf("[%.3fms] [rows: %v] %s | %v",
			float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
		return
	}

	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows: %v] %s",
		float64(elapsed.Nanoseconds())/1e6, rows, sql)
}
