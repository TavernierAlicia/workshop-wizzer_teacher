package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.NoRoute(func(c *gin.Context) {
		c.HTML(200, "error.html", nil)
	})

	// No more actions required pages
	router.GET("/", displayPage)
	router.GET("/index", displayPage)
	router.GET("/legal", displayPage)
	router.GET("/about", displayPage)

	// account stuff
	router.GET("/connect", displayPage)
	router.POST("/connect", connect)
	router.GET("/subscribe", displayPage)
	router.POST("/subscribe", subscribtion)
	router.GET("/forgotten-pwd", displayPage)
	router.POST("/forgotten-pwd", sendResetPWD)
	router.GET("/new-pwd/", getResetPage)
	router.POST("/new-pwd/", ResetPWD)

	// distinct views

	// all
	router.GET("/board/exercices")
	router.GET("/board/rank")
	router.GET("/board/params")

	// students
	router.GET("/board/histo")

	// alumn + profs
	router.GET("/board/overview")
	router.GET("/board/student/:id")

	// misc
	router.GET("/subInfos", showSubInfos)

	// student view
	// router.GET("/students/:id/overview", studentOverview)
	// router.GET("/students/:id/ranking", studentRanking)

	// test repo github
	// router.GET("/:github-account/:exerciseName")

	// misc

	router.Run(":9999")
}
