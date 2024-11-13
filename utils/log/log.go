package log

import (
	"github.com/sirupsen/logrus"
)

// InitLogger 设置日志格式和级别
func InitLogger(debug bool) {
	// 设置日志格式和级别
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

