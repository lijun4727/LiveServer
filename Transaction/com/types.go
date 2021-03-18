package com

import "github.com/gin-gonic/gin"

type ApiRoute interface {
	Init()
	Route(e *gin.Engine)
}
