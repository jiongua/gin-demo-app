package server

import (
	glog "gin_demo/interal/log"
	"github.com/gin-gonic/gin"
)

var log = glog.Log

func Start() {
	router := gin.New()
	router.Use(Logger())
	registerRoutes(router)
	router.Run()
}