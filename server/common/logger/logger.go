package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"tooladmin/server/common/config"

	"strings"
)

const (
	LevelDebug int = 0
	LevelInfo  int = 1
	LevelWarn  int = 2
	LevelError int = 3
)

var innerLogger *FileLogger

func InitLog() {
	innerLogger = NewFileLogger(config.GetConfig().Log)
}
func GetLoggerWriter() io.Writer {
	return innerLogger.logger.Writer()
}
func Debug(ctx context.Context, args ...interface{}) {
	if innerLogger.logLevel <= LevelDebug {
		innerLogger.Debug(ctx, args...)
	}
}

func Info(ctx context.Context, args ...interface{}) {
	if innerLogger.logLevel <= LevelInfo {
		innerLogger.Info(ctx, args...)
	}
}
func Warn(ctx context.Context, args ...interface{}) {
	if innerLogger.logLevel <= LevelWarn {

		innerLogger.Warn(ctx, args...)
	}
}
func Error(ctx context.Context, args ...interface{}) {
	if innerLogger.logLevel <= LevelError {
		innerLogger.Error(ctx, args...)
	}
}
func Fatal(ctx context.Context, args ...interface{}) {
	innerLogger.Fatal(ctx, args...)
	os.Exit(1)
}

func NewFileLogger(logConfig config.LogConfig) *FileLogger {
	err := os.MkdirAll(filepath.Dir(logConfig.Path), 0755)
	if err != nil {
		fmt.Printf("can not create log dir: %v,%v", filepath.Dir(logConfig.Path), err)
	}
	// 创建或打开日志文件，使用追加模式，权限0644
	file, err := os.OpenFile(logConfig.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("logConfig.Path", logConfig.Path)
		fmt.Printf("can not open log file: %v,%v", logConfig.Path, err)
	}
	resultLogger := &FileLogger{
		logger:   log.New(file, "", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.LUTC),
		logLevel: logConfig.Level,
	}
	return resultLogger
}

type FileLogger struct {
	logger   *log.Logger
	logLevel int
}

func (l *FileLogger) Debug(ctx context.Context, args ...interface{}) {
	format := make([]string, len(args))
	for i := range format {
		format[i] = "%+v"
	}
	format = append([]string{"[Debug] " + parseTraceID(ctx)}, format...)
	l.logger.Output(3, fmt.Sprintf(strings.Join(format, " "), args...))
}

func (l *FileLogger) Info(ctx context.Context, args ...interface{}) {
	format := make([]string, len(args))
	for i := range format {
		format[i] = "%+v"
	}
	format = append([]string{"[Info] " + parseTraceID(ctx)}, format...)
	l.logger.Output(3, fmt.Sprintf(strings.Join(format, " "), args...))

}
func (l *FileLogger) Warn(ctx context.Context, args ...interface{}) {
	format := make([]string, len(args))
	for i := range format {
		format[i] = "%+v"
	}
	format = append([]string{"[warn] " + parseTraceID(ctx)}, format...)
	l.logger.Output(3, fmt.Sprintf(strings.Join(format, " "), args...))

}

func (l *FileLogger) Error(ctx context.Context, args ...interface{}) {
	format := make([]string, len(args))
	for i := range format {
		format[i] = "%+v"
	}
	format = append([]string{"[Error] " + parseTraceID(ctx)}, format...)
	l.logger.Output(3, fmt.Sprintf(strings.Join(format, " "), args...))
}
func (l *FileLogger) Fatal(ctx context.Context, args ...interface{}) {
	format := make([]string, len(args))
	for i := range format {
		format[i] = "%+v"
	}
	format = append([]string{"[Fatal] " + parseTraceID(ctx)}, format...)
	l.logger.Output(3, fmt.Sprintf(strings.Join(format, " "), args...))
}

func parseTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		return "-"
	}
	return "[trace_id:" + traceID + "]"
}

// logWriter 自定义日志输出
type logWriter struct {
	writer io.Writer
}

func GetLogWriter(path string) *logWriter {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		fmt.Printf("can not create log dir: %v,%v", filepath.Dir(path), err)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("can not open log file: %v,%v", path, err)
	}
	return &logWriter{writer: file}
}

// Printf 实现logger.Writer接口
func (w *logWriter) Printf(format string, args ...interface{}) {
	format = "[mysql] " + time.Now().UTC().Format("2006-01-02 15:04:05") + " " + format + "\n"
	w.writer.Write([]byte(fmt.Sprintf(format, args...)))
}

func GetChannelMockLogger() *FileLogger {
	return NewFileLogger(config.LogConfig{
		Path:  "./logs/channel_mock.log",
		Level: LevelDebug,
	})
}
