package main

import (
	//	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	cmap "github.com/streamrail/concurrent-map"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles("frontend/dist/index.html")

	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", true)))

	tstate := TestState{MachineInfo: cmap.New()}
	go tstate.TimerCleaner()

	router.GET("/test/state/_data", tstate.StateData)
	// HeartBeat
	router.POST("/test/state/heartbeat", tstate.HeartBeat)
	router.GET("/test/state/heartbeat", tstate.GetHeartBeatInfo)

	router.POST("/test/state/guisummary", tstate.SetGUISummary)
	router.GET("/test/state/guisummary", tstate.GetGUISummary)
	router.DELETE("/test/state/guisummary/:testname_version", tstate.DelGUISummary)

	router.POST("/test/state/gui-machine-msg", tstate.SetGUIMachineMsg)
	router.GET("/test/state/gui-machine-msg", tstate.GetGUIMachineMsg)
	router.GET("/test/state/gui-machine-status", tstate.GetStatus)

	router.POST("/test/fuzz/chart-data", tstate.FuzzChartCreate)

	router.GET("/tools/xmlviewer", XMLViewer)

	FH := FilesHander{}
	router.GET("/files/getdata", FH.GetData)
	router.POST("/files/upload", FH.Upload)
	router.POST("/files/_update_db_info", FH.UpdateDBInfo)

	CIBD := CIBuilderHander{}
	router.GET("/ci/lastbuild", CIBD.GetLastBuild)
	router.POST("/ci", CIBD.CI)
	router.POST("/cistop", CIBD.CIStop)
	router.GET("/ci/_data", CIBD.CIState)

	router.GET("/", Index)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.Use(cors.Default())

	router.Run(":80")
}
