package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func disconnect(c *gin.Context) {

	session := sessions.Default(c)
	token := fmt.Sprintf("%v", session.Get("token"))

	id, err := checkToken(token)
	if err != nil || id == 0 {
		// error
		return
	}

	err = deleteToken(token)

	if err != nil {
		return
	}

	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/connect")
}

func removeExo(c *gin.Context) {
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

	// now insert in db
	err = deleteExo(infos.Id, exo_id)

	if err != nil {
		// html w err
		c.Redirect(http.StatusFound, "/board/exercices/del/true/ko")
		return
	}
	// success html
	c.Redirect(http.StatusFound, "/board/exercices/del/true/ok")
}

// TODO: delete account on demand
func deleteAccount(c *gin.Context) {

	c.Request.ParseForm()
	u_id, err := strconv.ParseInt(strings.Join(c.Request.PostForm["id"], " "), 10, 64)

	if err != nil {
		errToken(c)
		return
	}

	token := strings.Join(c.Request.PostForm["token"], " ")
	pwd := encodePWD(strings.Join(c.Request.PostForm["pwd"], " "))

	id, err := checkToken(token)
	if err != nil || id == 0 {
		errToken(c)
		return
	}

	if pwd == "" {
		c.HTML(200, "ask-delete.html", map[string]interface{}{"t": token, "send": 1, "ok": 0, "id": id})
		return
	}

	if u_id != id {
		errToken(c)
		return
	}

	infos, err := getUserInfos(token)

	if err != nil {
		errToken(c)
		return
	}

	// verify password too

	_, err = getConnected(infos.Mail, pwd)

	if err != nil {
		c.HTML(200, "ask-delete.html", map[string]interface{}{"t": token, "send": 1, "ok": 0, "id": id})
		return
	}

	// database
	err = deleteAllUserData(id)

	if err != nil {
		c.HTML(200, "ask-delete.html", map[string]interface{}{"t": token, "send": 1, "ok": 0, "id": id})
		return
	}

	message := ` 
	<p> Conformément à votre demande, votre compte ainsi que toutes les données accociées à celui-ci ont été supprimés. </p> 
	<p> Wizzer Teacher vous souhaite une bonne continuation! </p>
	`
	_ = confirmDelete(infos.Mail, message)

	c.HTML(200, "ask-delete.html", map[string]interface{}{"t": token, "send": 1, "ok": 1, "id": id})

}
