package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var subForm Sub
var err error

func subscribtion(c *gin.Context) {

	// Set data
	c.Request.ParseForm()

	if strings.Join(c.Request.PostForm["legal"], " ") != "on" {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	subForm.AccountType = strings.Join(c.Request.PostForm["account"], " ")
	subForm.Name = strings.Join(c.Request.PostForm["name"], " ")
	subForm.Surname = strings.Join(c.Request.PostForm["surname"], " ")
	subForm.Mail = strings.Join(c.Request.PostForm["mail"], " ")
	subForm.Repo = strings.Join(c.Request.PostForm["repo"], " ")
	subForm.Campus = strings.Join(c.Request.PostForm["campus"], " ")
	subForm.Studies = strings.Join(c.Request.PostForm["formation"], " ")
	subForm.Matter = strings.Join(c.Request.PostForm["matiere"], " ")
	subForm.Pwd = strings.Join(c.Request.PostForm["pwd"], " ")
	subForm.PwdConfirm = strings.Join(c.Request.PostForm["pwd-confirm"], " ")

	// verify data

	if (subForm.AccountType == "" && len(subForm.AccountType) > 250) ||
		(subForm.Name == "" || len(subForm.Name) > 250) ||
		(subForm.Surname == "" || len(subForm.Surname) > 250) ||
		(subForm.Mail == "" || len(subForm.Mail) > 250) ||
		(subForm.Repo == "" && subForm.AccountType == "student") ||
		(len(subForm.Repo) > 250) ||
		(subForm.Campus == "") ||
		(subForm.Studies == "" && subForm.AccountType == "student") ||
		(subForm.Matter == "" && subForm.AccountType != "student") ||
		(subForm.Pwd == "" || len(subForm.Pwd) > 250) ||
		(subForm.PwdConfirm == "" || len(subForm.PwdConfirm) > 250) {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// radio
	if stringInSlice(subForm.AccountType, []string{"student", "prof", "alum"}) {
		// empty git repo if is not a student
		if subForm.AccountType != "student" {
			subForm.Repo = ""
			subForm.Studies = ""
		} else {
			subForm.Matter = ""
		}
	} else {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// mail
	mailreg, _ := regexp.Compile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])")
	if !mailreg.MatchString(subForm.Mail) {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// get lists to verify data
	campuslist, err := getSchools()
	studieslist, err := getStudies()
	matterslist, err := getMatters()

	// select inputs
	if stringInSlice(subForm.Campus, campuslist) {

		if subForm.AccountType == "student" {
			if !stringInSlice(subForm.Studies, studieslist) {
				c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
				return
			}
		} else {
			if !stringInSlice(subForm.Matter, matterslist) {
				c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
				return
			}
		}
	} else {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// pwd
	if subForm.Pwd != subForm.PwdConfirm {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// record data
	err = RecordUser(subForm)

	if err != nil {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 1})

}

func connect(c *gin.Context) {

	c.Request.ParseForm()

	mail := strings.Join(c.Request.PostForm["connect-id"], " ")
	pwd := strings.Join(c.Request.PostForm["connect-pwd"], " ")

	if mail == "" || pwd == "" {
		c.HTML(200, "connect.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	pwd = encodePWD(pwd)

	token, err := getConnected(mail, pwd)

	if err != nil {
		c.HTML(200, "connect.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	infos, err := getUserInfos(token)

	if err != nil {
		c.HTML(200, "connect.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge: 86400 * 30,
	})

	session.Set("token", token)
	session.Set("type", infos.Type)
	session.Set("campus_id", infos.CampusID)
	session.Set("matter_id", infos.MatterID)
	session.Set("studies_id", infos.StudiesID)
	session.Set("name", infos.Name)
	session.Set("surname", infos.Surname)
	session.Set("pic", infos.Pic)
	session.Save()

	// now use session then display the good dashboard

	if infos.Type == "student" || infos.Type == "alum" || infos.Type == "prof" {
		c.Redirect(http.StatusFound, "/board/exercices")
	} else {
		c.HTML(200, "connect.html", map[string]interface{}{"send": 1, "ok": 0})
	}

}

func sendResetPWD(c *gin.Context) {
	c.Request.ParseForm()

	mail := strings.Join(c.Request.PostForm["connect-id"], " ")

	if mail == "" || getMail(mail) == "" {
		// error
		c.HTML(200, "forgotten-pwd.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// generate new token to reset mail
	newtoken := tokenGenerator()

	link := viper.GetString("links.host") + "new-pwd/?token=" + newtoken

	// insert new token in db
	err = insertToken(newtoken, mail)
	if err != nil {
		c.HTML(200, "forgotten-pwd.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// send mail
	err = SendResetMail(mail, link)

	// if error
	if err != nil {
		// delete token from db
		_ = deleteToken(mail)
		c.HTML(200, "forgotten-pwd.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// okay
	c.HTML(200, "forgotten-pwd.html", map[string]interface{}{"send": 1, "ok": 1})

}

func ResetPWD(c *gin.Context) {
	c.Request.ParseForm()

	token := c.Query("token")
	pwd := strings.Join(c.Request.PostForm["pwd"], " ")
	confirm_pwd := strings.Join(c.Request.PostForm["pwd-confirm"], " ")

	// check pwd
	if (pwd == "" || confirm_pwd == "") || pwd != confirm_pwd {
		c.HTML(200, "new-pwd.html", map[string]interface{}{"send": 1, "ok": 0, "path": c.Query("token")})
		return
	}

	id, _ := checkToken(token)

	if id == 0 {
		c.HTML(200, "new-pwd.html", map[string]interface{}{"send": 1, "ok": 0, "path": c.Query("token")})
		return
	}

	err = updatePWD(pwd, id)
	if err != nil {
		c.HTML(200, "new-pwd.html", map[string]interface{}{"send": 1, "ok": 0, "path": c.Query("token")})
		return
	}

	err = deleteToken(token)

	c.HTML(200, "new-pwd.html", map[string]interface{}{"send": 1, "ok": 1, "path": ""})
}

func addExo(c *gin.Context) {
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
		c.Redirect(http.StatusFound, "/board/exercices/add/true/ko")
		return
	}

	// now insert in db
	err = postExo(exoName, repo, exoDate, description, level, exoMatter, exoLang, bar, infos.Id)

	if err != nil {
		// html w err
		c.Redirect(http.StatusFound, "/board/exercices/add/true/ko")
		return
	}
	// success html
	c.Redirect(http.StatusFound, "/board/exercices/add/true/ok")
}

func recordGrade(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")

	_, err := checkToken(token)
	if err != nil {
		c.JSON(401, nil)
		return
	}

	var empty NewGrade
	var exercice NewGrade

	c.BindJSON(&exercice)
	if exercice == empty ||
		exercice.ExerciceID == 0 ||
		exercice.StudentID == 0 ||
		exercice.Score == 0 {
		c.JSON(400, nil)
		return
	}

	err = insertRendu(exercice)
	if err != nil {
		c.JSON(500, nil)
		return
	}

	c.JSON(201, nil)

}

func resetBotToken(c *gin.Context) {
	data := GetSessionData(sessions.Default(c))

	id, err := checkToken(data.Token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	infos, err := getUserInfos(data.Token)

	if err != nil || data.Atype != "prof" {
		errToken(c)
		return
	}

	_, err = updateBotToken(infos.Id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusFound, "/board/params")

}
