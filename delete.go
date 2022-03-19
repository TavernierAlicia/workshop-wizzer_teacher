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
