package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

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
