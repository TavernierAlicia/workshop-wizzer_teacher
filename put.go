package main

import (
	"strings"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func recordParams(c *gin.Context) {

	c.Request.ParseForm()
	pic := strings.Join(c.Request.PostForm["pic"], " ")
	repo := strings.Join(c.Request.PostForm["url"], " ")
	campus := strings.Join(c.Request.PostForm["campus"], " ")
	studies := strings.Join(c.Request.PostForm["formation"], " ")
	matter := strings.Join(c.Request.PostForm["matter"], " ")

	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	campuslist, err := getSchools()
	studieslist, err := getStudies()
	matterslist, err := getMatters()

	infos, _ := getUserInfos(data.Token)

	// select inputs
	if !stringInSlice(campus, campuslist) {
		c.HTML(200, "parameters.html", map[string]interface{}{"send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
		return
	}

	if data.Atype == "student" {
		if !stringInSlice(studies, studieslist) || repo == "" {
			c.HTML(200, "parameters.html", map[string]interface{}{"send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	} else {
		if !stringInSlice(matter, matterslist) {
			c.HTML(200, "parameters.html", map[string]interface{}{"send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	}

	// now update in db
	if pic == "" {
		pic = viper.GetString("default.default_pic")
	}
	err = updateParams(infos.Id, pic, repo, campus, studies, matter)

	if err != nil {
		c.HTML(200, "parameters.html", map[string]interface{}{"send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
		return
	}

	// TODO !!!! insert imgs in server

	infos, _ = getUserInfos(data.Token)
	sessions.Default(c).Set("pic", infos.Pic)
	sessions.Default(c).Save()

	c.HTML(200, "parameters.html", map[string]interface{}{"send": 1, "ok": 1, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})

}
