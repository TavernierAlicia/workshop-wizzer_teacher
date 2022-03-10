package main

import (
	"github.com/gin-gonic/gin"
)

func displayPage(c *gin.Context) {

	page := c.Request.URL.Path[1:]
	if page == "" {
		page = "index"
	}
	if page == "subscribe" {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 0, "ok": 0})
	} else {
		c.HTML(200, page+".html", nil)
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
