package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func displayPage(c *gin.Context) {

	page := c.Request.URL.Path[1:]
	if page == "" {
		page = "index"
	}

	switch page {
	case "subscribe":
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 0, "ok": 0})
	case "connect":
		c.HTML(200, "connect.html", map[string]interface{}{"send": 0, "ok": 0})
	case "forgotten-pwd":
		c.HTML(200, "forgotten-pwd.html", map[string]interface{}{"send": 0, "ok": 0})

	default:
		c.HTML(200, page+".html", nil)
	}
}

func getResetPage(c *gin.Context) {
	// check token before
	id, _ := checkToken(c.Query("token"))
	if id != 0 {
		c.HTML(200, "new-pwd.html", map[string]interface{}{"send": 0, "ok": 0, "path": c.Query("token")})
	} else {
		c.AbortWithError(http.StatusBadRequest, err)
	}
}

func showSubInfos(c *gin.Context) {

	schools, err := getSchools()
	if err != nil {
		c.JSON(500, nil)
		return
	}
	formations, err := getStudies()
	if err != nil {
		c.JSON(500, nil)
		return
	}
	matters, err := getMatters()
	if err != nil {
		c.JSON(500, nil)
		return
	}

	infos := &SubFormInfos{
		Schools:    schools,
		Formations: formations,
		Matters:    matters,
	}

	c.JSON(200, infos)
}
