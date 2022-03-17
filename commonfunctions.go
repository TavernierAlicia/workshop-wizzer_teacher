package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var logger *zap.Logger

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func encodePWD(pwd string) (string, error) {

	sum := sha256.Sum256([]byte(pwd))
	pass := hex.EncodeToString(sum[:])
	fmt.Println(pass)

	if err != nil {
		fmt.Println("error while converting pwd")
		return "", err
	} else {
		return pass, err
	}
}

func tokenGenerator() string {
	token := uuid.New()
	return token.String()
}

// print errors
func printErr(desc string, nomFunc string, err error) {
	logger, _ = zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Error("Cannot "+desc, zap.String("Func", nomFunc), zap.Error(err))
	}
}

func errToken(c *gin.Context) {
	c.HTML(403, "unauthorized.html", nil)
}

func GetSessionData(session sessions.Session) (data SessionInfos) {
	data.Token = fmt.Sprintf("%v", session.Get("token"))
	data.Atype = fmt.Sprintf("%v", session.Get("type"))
	data.Name = fmt.Sprintf("%v", session.Get("name"))
	data.Surname = fmt.Sprintf("%v", session.Get("surname"))
	data.Campus_id = fmt.Sprintf("%v", session.Get("campus_id"))
	data.Matter_id = fmt.Sprintf("%v", session.Get("matter_id"))
	data.Studies_id = fmt.Sprintf("%v", session.Get("studies_id"))
	return data
}
