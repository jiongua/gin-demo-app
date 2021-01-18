package entity

import (
	"fmt"
	"gorm.io/gorm/logger"
	golog "log"
	"os"
	"path/filepath"
	"time"
)

var GormLogFile *os.File

func init()  {
	rootPath, _ := filepath.Abs(filepath.Dir("../../"))
	path := filepath.Join(rootPath, "logs", "gorm.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		GormLogFile = file
		fmt.Println("set log file ok")
	}else {
		GormLogFile = os.Stdout
		log.Errorf("set log file error: %v", err)
		fmt.Println("set log file error")
	}
	logger.Default.LogMode(logger.Info)
}

var GormLogger = logger.New(
	golog.New(GormLogFile, "\r\n", golog.LstdFlags), // io writer
	logger.Config{
		SlowThreshold: time.Second,   // 慢 SQL 阈值
		LogLevel:      logger.Silent, // Log level
		Colorful:      false,         // 禁用彩色打印
	},
)
