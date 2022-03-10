package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// load webdev stuff
	router.LoadHTMLGlob("templates/*")
	router.Static("/js", "./js")
	router.Static("./css", "./css")

	//to include images
	router.Static("/pictures", "./pictures")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(200, "error.html", nil)
	})

	// No more actions required pages
	router.GET("/", displayPage)
	router.GET("/index", displayPage)
	router.GET("/legal", displayPage)
	router.GET("/team", displayPage)
	router.GET("/about", displayPage)
	router.GET("/connect", displayPage)
	router.GET("/subscribe", displayPage)

	router.POST("/subscribe", subscribtion)

	// getting infos for sub
	router.GET("/subInfos", showSubInfos)

	// student view
	// router.GET("/students/:id/overview", studentOverview)
	// router.GET("/students/:id/ranking", studentRanking)

	// test repo github
	// router.GET("/:github-account/:exerciseName")

	// misc

	router.Run(":9999")
}
