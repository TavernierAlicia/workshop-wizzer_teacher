package main

import (
	"fmt"
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

func showLevels(c *gin.Context) {
	levels, err := getLevels()
	if err != nil {
		c.JSON(500, nil)
		return
	}

	c.JSON(200, levels)
}

func getExos(c *gin.Context) {

	// verify type & token
	session := sessions.Default(c)
	token := fmt.Sprintf("%v", session.Get("token"))
	aType := fmt.Sprintf("%v", session.Get("type"))
	name := fmt.Sprintf("%v", session.Get("name"))
	surname := fmt.Sprintf("%v", session.Get("surname"))
	campusId := fmt.Sprintf("%v", session.Get("campus_id"))
	matterId := fmt.Sprintf("%v", session.Get("matter_id"))
	studiesId := fmt.Sprintf("%v", session.Get("studies_id"))

	id, err := checkToken(token)
	if err != nil || id == 0 {
		// error
		return
	}

	exos, _ := getExercices(id, aType, studiesId, campusId, matterId)

	size := len(exos)

	if aType == "student" {
		// get day exos
		score := getScore(id)
		c.HTML(200, "board.html", map[string]interface{}{"name": name, "surname": surname, "score": score, "size": size, "student": 1, "exos": exos})
	} else {

		first := strconv.Itoa(time.Now().Year()) + "/06" + "/01"
		last := strconv.Itoa(time.Now().Year()) + "/10" + "/01"

		if len(exos) != 0 {
			first = exos[0].Due
			last = exos[len(exos)-1].Due
		}

		if aType == "prof" {
			c.HTML(200, "board.html", map[string]interface{}{"name": name, "surname": surname, "size": size, "student": 0, "last": last, "first": first, "exos": exos})
		} else if aType == "alum" {
			c.HTML(200, "board.html", map[string]interface{}{"name": name, "surname": surname, "size": size, "student": 9, "last": last, "first": first, "exos": exos})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}

}
