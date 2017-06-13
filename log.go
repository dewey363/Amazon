package main

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func InitLog(logFile string) *logrus.Logger {
	// Create a new instance of the logger. You can have any number of instances.
	fileHook := NewHook()
	//	logrus.AddHook(fileHook)
	var log = logrus.New()
	log.Hooks.Add(fileHook)
	saveDir := "log/"
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+saveDir, os.ModePerm) //生成多级目录
	if err != nil {
		log.Info(err)
		//		return "", errors.New("输出到文件错误")
	}
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(saveDir+logFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	return log
}
