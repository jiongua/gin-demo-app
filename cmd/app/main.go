package main

import (
	_ "gin_demo/interal/crawl"
	_ "gin_demo/interal/entity"
	"gin_demo/interal/server"
)

func main() {
	server.Start()
}
