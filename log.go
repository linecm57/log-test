package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"runtime"
	"net"
	"path/filepath"
)

// CustomFormat 继承 logrus.JSONFormatter，用于自定义日志格式
type CustomFormat struct {
	logrus.JSONFormatter
}

// Format 实现自定义格式
func (f *CustomFormat) Format(entry *logrus.Entry) ([]byte, error) {
	// 检查并替换 msg 为 fields.msg 的内容
	if msgField, ok := entry.Data["msg"]; ok {
		entry.Message = fmt.Sprintf("%v", msgField)
		delete(entry.Data, "msg") // 从自定义字段中移除 msg
	}
	if msgField, ok := entry.Data["time"]; ok {
		entry.Message = fmt.Sprintf("%v", msgField)
		delete(entry.Data, "time") // 从自定义字段中移除 msg
	}

	// 使用 JSONFormatter 来格式化其余的内容
	return f.JSONFormatter.Format(entry)
}

// 初始化 logrus 日志
func initLogger() *logrus.Logger {
	logger := logrus.New()

	// 设置输出到文件或标准输出
	logger.Out = os.Stdout

	// 设置日志格式为自定义格式
	logger.SetFormatter(&CustomFormat{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: time.RFC3339, // 使用 RFC3339 格式的时间戳
		},
	})

	return logger
}

func main() {
	// 创建 logrus 实例
	logger := initLogger()

	// 记录日志并添加自定义字段
	logger.WithFields(logrus.Fields{
		"ip":       getOutboundIP(),
		"path":     getFileInfo(),
		"severity": "info",
		"msg":      "heelp", // 这里的 msg 会被替换到日志的核心字段
		"time":     time.Now().Format(time.RFC3339),
	}).Info()
}

// 获取文件名和行号
func getFileInfo() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// 获取本机IP地址
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
