package main

import (
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
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
		session := sessions.Default(c)
		token := session.Get("token")

		if reflect.TypeOf(token) != nil {
			c.Redirect(http.StatusFound, "/board/exercices")
		}

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

func showLvlLang(c *gin.Context) {

	languages, err := getLanguages()
	if err != nil {
		c.JSON(500, nil)
		return
	}

	levels, err := getLevels()
	if err != nil {
		c.JSON(500, nil)
		return
	}

	infos := &SubFormInfos{
		Levels:    levels,
		Languages: languages,
	}
	c.JSON(200, infos)
}

func getExos(c *gin.Context) {

	// verify type & token
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	if data.Atype == "student" {
		params := exoSearch{}
		exos, _ := getExercices(id, data.Atype, data.Studies_id, data.Campus_id, data.Matter_id, params)
		size := len(exos)
		score := getScore(id)
		c.HTML(200, "board.html", map[string]interface{}{"name": data.Name, "surname": data.Surname, "score": score, "size": size, "student": 1, "exos": exos})
	} else {

		// getting possible search criterions
		params := exoSearch{
			Name:     c.Query("exo-name"),
			Level:    c.Query("exo-level"),
			Date:     c.Query("date"),
			Language: c.Query("exo-language"),
		}

		// get exos
		exos, _ := getExercices(id, data.Atype, data.Studies_id, data.Campus_id, data.Matter_id, params)
		size := len(exos)
		first := strconv.Itoa(time.Now().Year()) + "/06" + "/01"
		last := strconv.Itoa(time.Now().Year()) + "/10" + "/01"

		if len(exos) != 0 {
			first = exos[0].Due
			last = exos[len(exos)-1].Due
		}

		if data.Atype == "prof" {
			c.HTML(200, "board.html", map[string]interface{}{"name": data.Name, "surname": data.Surname, "size": size, "student": 0, "last": last, "first": first, "exos": exos})
		} else if data.Atype == "alum" {
			c.HTML(200, "board.html", map[string]interface{}{"name": data.Name, "surname": data.Surname, "size": size, "student": 9, "last": last, "first": first, "exos": exos})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}

}
