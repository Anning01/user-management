package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// Init 初始化日志
func Init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 记录信息日志
func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

// Infof 格式化记录信息日志
func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

// Error 记录错误日志
func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

// Errorf 格式化记录错误日志
func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

// Fatal 记录致命错误并退出
func Fatal(v ...interface{}) {
	errorLogger.Fatal(v...)
}

// Fatalf 格式化记录致命错误并退出
func Fatalf(format string, v ...interface{}) {
	errorLogger.Fatalf(format, v...)
}
