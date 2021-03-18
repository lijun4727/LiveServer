package main

// 周期原因，代码粗糙,记得调整

import (
	"LiveServer/common"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	hub := newHub()
	go hub.run()

	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()
	e.GET("/socket.io", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})
	e.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./static/index.html")
	})
	e.Static("/static/", "./static")
	//endless.ListenAndServeTLS(common.SignalServicePort, "ssl/server.crt", "ssl/server.key", e)
	endless.ListenAndServe(common.SignalPort, e)
}
