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
	// actually DELETE but html forms makes me sad
	router.GET("/board/disconnect", disconnect)

	router.GET("/board/exercices", getExos)
	router.GET("/board/rank")
	router.GET("/board/params")

	// students
	router.GET("/board/histo")

	// alumn + profs
	router.GET("/board/overview")
	router.GET("/board/student/:id")

	// misc
	router.GET("/subInfos", showSubInfos)
	router.GET("/getLvlLang", showLvlLang)
	// router.POST("/results/exercice_id/student_id/", recordGrade)

	router.Run(":9999")
}
