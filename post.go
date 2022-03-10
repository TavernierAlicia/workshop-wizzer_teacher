package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var subForm Sub
var err error

func subscribtion(c *gin.Context) {

	// Set data
	c.Request.ParseForm()

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

	fmt.Println(subForm)
	// verify data

	// null
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
		fmt.Println("OK")
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
	if subForm.Pwd == subForm.PwdConfirm {
		fmt.Println("OK")

		// sha pwd ready for storage
		subForm.Pwd, err = encodePWD(subForm.Pwd)
		if err != nil {
			c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
			return
		}

	} else {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	// record data
	fmt.Println("ready to store data")
	err = RecordUser(subForm)

	if err != nil {
		c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 0})
		return
	}

	c.HTML(200, "subscribe.html", map[string]interface{}{"send": 1, "ok": 1})

}
