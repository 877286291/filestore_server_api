package main

import (
	"github.com/aurora/Filestore-server/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	routers.InitRouters(r)
	_ = r.Run()
}
