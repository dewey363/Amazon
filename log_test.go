package main

import (
	//	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

//测试单个插入
func TestInitLog(t *testing.T) {
	log := InitLog("test.log")
	log.WithFields(logrus.Fields{
		"animal": "嘻嘻嘻",
		"size":   30,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"xxx":  "xxx",
		"size": 333,
	}).Warnf("有行号%d", 12)
}
