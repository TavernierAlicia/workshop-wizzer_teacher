package main

import (
	"fmt"
	"net/http"

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
	fmt.Println("DELETE")
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
