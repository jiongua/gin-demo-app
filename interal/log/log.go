package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var Log *logrus.Logger


func init() {
	Log = &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    &logrus.TextFormatter{},
		Level:        logrus.DebugLevel,
	}
	//path := FileName()
	//file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err == nil {
	//	Log.Out = file
	//} else {
	//	Log.Infof("Failed to log to file: %v", err)
	//}
}

func FileName() string {
	rootPath, _ := filepath.Abs(filepath.Dir("."))
	return filepath.Join(rootPath, "logs", "api.log")
}

