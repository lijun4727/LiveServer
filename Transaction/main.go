package main

import (
	"LiveServer/Transaction/com"
	"LiveServer/Transaction/db_connect"
	"LiveServer/Transaction/route"
	"LiveServer/common"
	"flag"
	"fmt"
	"log"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	var test = flag.String("test", "test_var", "test var ")
	flag.Parse()

	err := db_connect.InitConnect()
	if err != nil {
		log.Fatalf("failed to connect mysql: %v", err)
		return
	}
	defer db_connect.Clean()

	e := gin.Default()
	e.GET("/test", func(c *gin.Context) {
		testStr := fmt.Sprintf("test=%s", *test)
		c.String(200, testStr)
	})

	com.AppendRoute(e, &route.AccountRoute{})
	com.AppendRoute(e, &route.ContactPersons{})
	endless.DefaultReadTimeOut = common.ApiReadTimeOut
	endless.DefaultWriteTimeOut = common.ApiWriteTimeOut
	endless.DefaultMaxHeaderBytes = 1 << 20
	endless.ListenAndServe(common.MediaTransactionPort, e)
}
