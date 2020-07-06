package main

import (

	//"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	return
	//	c.String(http.StatusOK, "xx QA WebSite")
}
