package server

import (
	glog "gin_demo/interal/log"
	"gin_demo/interal/task"
	"github.com/gin-gonic/gin"
)

var log = glog.Log

func Start() {
	router := gin.New()
	router.Use(Logger())
	registerRoutes(router)
	task.StartAnswerNotifyWorker()
	task.StartVoteAnswerNotifyWorker()
	router.Run()
}