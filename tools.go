package main

import (

	//"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func XMLViewer(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://codebeautify.org/xmlviewer")

}
