package main

import (
	"net/http"
	"strings"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func recordParams(c *gin.Context) {

	pic := ""
	ext := ""
	repo := ""
	studies := ""
	matter := ""
	campus := ""

	data := GetSessionData(sessions.Default(c))
	infos, _ := getUserInfos(data.Token)

	campuslist, _ := getSchools()
	studieslist, _ := getStudies()
	matterslist, _ := getMatters()

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	bt, _ := getBotToken(id)

	result, _ := c.MultipartForm()

	if len(result.Value["url"]) == 0 {
		repo = ""
		if infos.Type == "student" {
			c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	} else {
		repo = result.Value["url"][0]
	}

	if len(result.Value["formation"]) == 0 {
		studies = ""
		if infos.Type == "student" {
			c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	} else {
		studies = result.Value["formation"][0]
	}

	if len(result.Value["matter"]) == 0 {
		matter = ""
		if infos.Type != "student" {
			c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	} else {
		matter = result.Value["matter"][0]
	}

	if len(result.Value["campus"]) == 0 {
		c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
		return
	} else {
		campus = result.Value["campus"][0]
	}

	// TODO: rename file, insert correct path in db
	file, err := c.FormFile("pic")
	if err != nil {
		printErr("get pic formfile", "recordParams", err)
	} else {
		if len(file.Filename) > 0 {
			ext = after(file.Filename, ".")

			pic = "saved/" + tokenGenerator() + "." + ext

			err = c.SaveUploadedFile(file, pic)

			if err != nil {
				printErr("save pic formfile", "recordParams", err)
			}
		}
	}

	// select inputs
	if !stringInSlice(campus, campuslist) {
		c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
		return
	}

	if data.Atype == "student" {
		if !stringInSlice(studies, studieslist) || repo == "" {
			c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	} else {
		if !stringInSlice(matter, matterslist) {
			c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
			return
		}
	}

	// now update in db
	if pic == "" {
		pic = viper.GetString("default.default_pic")
	}
	err = updateParams(infos.Id, pic, repo, campus, studies, matter)

	if err != nil {
		c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 0, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})
		return
	}

	infos, _ = getUserInfos(data.Token)
	sessions.Default(c).Set("pic", infos.Pic)
	sessions.Default(c).Save()

	c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": bt, "send": 1, "ok": 1, "campus": campuslist, "matter": matterslist, "studies": studieslist, "infos": infos})

}

func editExo(c *gin.Context) {
	c.Request.ParseForm()
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil || data.Atype != "prof" {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	exo_id := c.Query("exo-id")

	description := strings.Join(c.Request.PostForm["desc"], " ")
	exoName := strings.Join(c.Request.PostForm["exo-name"], " ")
	bar := strings.Join(c.Request.PostForm["bar"], " ")
	exoDate := strings.Join(c.Request.PostForm["exo-date"], " ")
	exoMatter := strings.Join(c.Request.PostForm["exo-matter"], " ")
	exoLang := strings.Join(c.Request.PostForm["exo-language"], " ")
	level := strings.Join(c.Request.PostForm["exo-level"], " ")
	repo := strings.Join(c.Request.PostForm["repo-path"], " ")

	matterslist, _ := getMatters()
	levelslist, _ := getLevels()
	languageslist, _ := getLanguages()

	if (description == "" || len(description) > 500) ||
		(exoName == "" || len(exoName) > 250) ||
		(bar == "" || len(bar) > 3) ||
		(exoDate == "" || len(exoDate) != 10) ||
		(!stringInSlice(exoMatter, matterslist)) ||
		(!stringInSlice(exoLang, languageslist)) ||
		(!stringInSlice(level, levelslist)) ||
		(repo == "" || len(exoDate) > 250) {
		// html w err
		c.Redirect(http.StatusFound, "/board/exercices/edit/true/ko")
		return
	}

	// now insert in db
	err = putExo(exoName, repo, exoDate, description, level, exoMatter, exoLang, bar, infos.Id, exo_id)

	if err != nil {
		// html w err
		c.Redirect(http.StatusFound, "/board/exercices/edit/true/ko")
		return
	}
	// success html
	c.Redirect(http.StatusFound, "/board/exercices/edit/true/ok")
}
