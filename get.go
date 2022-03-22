package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	languagesList, _ := getLanguages()
	mattersList, _ := getMatters()
	levelsList, _ := getLevels()

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		tokenMismatch(c)
		return
	}

	infos, _ := getUserInfos(data.Token)

	if data.Atype == "student" {
		params := exoSearch{}
		exos, _ := getExercices(id, data.Atype, data.Studies_id, data.Campus_id, data.Matter_id, params)
		size := len(exos)
		score := getScore(id)

		c.HTML(200, "board.html", map[string]interface{}{
			"size":          size,
			"student":       1,
			"exos":          exos,
			"score":         score,
			"infos":         infos,
			"mattersList":   mattersList,
			"levelsList":    levelsList,
			"languagesList": languagesList,
			"is_delete":     0,
			"is_edit":       0,
			"id_add":        0,
			"is_success":    0,
			"is_send":       0,
		})
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
			is_add := 0
			is_delete := 0
			is_edit := 0
			is_success := 0
			is_send := 0
			exo_details := Exos{}

			switch c.Param("action") {
			case "add":
				is_add = 1

			case "edit":
				is_edit = 1
				if c.Query("exo-id") != "" {
					exo_details, err = getExoDetails(c.Query("exo-id"), id)
					if err != nil {
						is_edit = 0
					}
				}

			case "del":
				is_delete = 1
				if c.Query("exo-id") != "" {
					exo_details, err = getExoDetails(c.Query("exo-id"), id)
					if err != nil {
						is_delete = 0
					}
				}

			default:
			}

			if c.Param("send") == "true" {
				is_send = 1
				if c.Param("result") == "ok" {
					is_success = 1
				}
			}

			c.HTML(200, "board.html", map[string]interface{}{
				"size":          size,
				"student":       0,
				"last":          last,
				"first":         first,
				"exos":          exos,
				"infos":         infos,
				"mattersList":   mattersList,
				"levelsList":    levelsList,
				"languagesList": languagesList,
				"is_delete":     is_delete,
				"is_edit":       is_edit,
				"is_add":        is_add,
				"is_success":    is_success,
				"is_send":       is_send,
				"exo_details":   exo_details,
			})
		} else if data.Atype == "alum" {
			c.HTML(200, "board.html", map[string]interface{}{
				"size":          size,
				"student":       9,
				"last":          last,
				"first":         first,
				"exos":          exos,
				"infos":         infos,
				"mattersList":   mattersList,
				"levelsList":    levelsList,
				"languagesList": languagesList,
				"is_delete":     0,
				"is_edit":       0,
				"id_add":        0,
				"is_success":    0,
				"is_send":       0,
			})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

}

func getParams(c *gin.Context) {
	// verify type & token
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	schools, _ := getSchools()
	formations, _ := getStudies()
	matters, _ := getMatters()

	botToken := ""
	if infos.Type == "prof" {
		botToken, err = getBotToken(infos.Id)

		if err != nil {
			botToken = "un problème est survenu, nous vous conseillons de rafraîchir votre botToken"
		}
	}

	c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "no", "botToken": botToken, "send": 0, "ok": 0, "campus": schools, "matter": matters, "studies": formations, "infos": infos})

}

func getRank(c *gin.Context) {
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// ok, fetch data
	students, err := getStudents(infos.StudiesID, infos.CampusID, infos.MatterID)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	student := 0
	var score int64
	if data.Atype == "student" {
		student = 1
		score = getScore(id)

	}

	fmt.Println(students)

	// display html
	c.HTML(200, "rank.html", map[string]interface{}{
		"infos": infos, "student": student, "score": score, "students": students,
	})

}

func getStudentHisto(c *gin.Context) {
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	student_id := c.Query("id")
	student := 0

	if data.Atype == "student" {
		student = 1
		if strconv.FormatInt(infos.Id, 10) != student_id {
			errToken(c)
			return
		}
	}

	// now get data
	score := getScore(id)
	studentScore, err := getStudentScoring(student_id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// return html
	c.HTML(200, "student.html", map[string]interface{}{
		"infos": infos, "student": student, "score": score, "studentScoring": studentScore,
	})
}

func getOverview(c *gin.Context) {
	data := GetSessionData(sessions.Default(c))

	// getting possible search criterions
	params := OverviewSearch{
		Level:   c.Query("exo-level"),
		Studies: c.Query("exo-studies"),
	}

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if data.Atype == "student" {
		errToken(c)
		return
	}

	// now get data
	studentScore, err := getAllStudentScoring(infos.CampusID, infos.MatterID, params)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	studiesList, _ := getStudies()
	levelsList, _ := getLevels()

	// return html
	c.HTML(200, "overview.html", map[string]interface{}{
		"levelsList": levelsList, "studiesList": studiesList, "infos": infos, "studentScoring": studentScore,
	})
}

// TODO: get data import to csv

// TODO: ask delete account sendmail
func askDeleteAccount(c *gin.Context) {

	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	deleteLink := viper.GetString("links.host") + "delete-account/view?id=" + strconv.FormatInt(infos.Id, 10) + "&t=" + data.Token

	err = sendDeleteMail(infos.Mail, deleteLink)

	botToken := ""

	if infos.Type == "prof" {
		botToken, err = getBotToken(infos.Id)

		if err != nil {
			botToken = "un problème est survenu, nous vous conseillons de rafraîchir votre botToken"
		}
	}

	schools, _ := getSchools()
	formations, _ := getStudies()
	matters, _ := getMatters()

	if err != nil {
		c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "failed", "botToken": botToken, "send": 0, "ok": 0, "campus": schools, "matter": matters, "studies": formations, "infos": infos})
	}

	c.HTML(200, "parameters.html", map[string]interface{}{"getData": "no", "deleteAccount": "success", "botToken": botToken, "send": 0, "ok": 0, "campus": schools, "matter": matters, "studies": formations, "infos": infos})
}

func DeleteView(c *gin.Context) {

	token := c.Query("t")
	paramid, err := strconv.ParseInt(c.Query("id"), 10, 64)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := checkToken(token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	if paramid != id {
		errToken(c)
		return
	}

	c.HTML(200, "ask-delete.html", map[string]interface{}{"t": token, "send": 0, "ok": 0, "id": id})
}
