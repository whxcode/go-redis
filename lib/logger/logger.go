package logger

import (
	"log"
	"os"

	"gorm.io/gorm/logger"
)

var l logger.Interface

func init() {
	l = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 关键：用 log.New 包装 os.Stdout
		logger.Config{})
}
