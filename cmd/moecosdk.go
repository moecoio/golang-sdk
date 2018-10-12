package main

import (
	"sdk"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	/*
	 * All logs redirected to stdout
	 */
	//logger.SetLevel(logrus.DebugLevel)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	MoecoSdk := sdk.NewMoecoSDK(
		"https://prod114.moeco.io:443",
		"API_KEY",
		"NODE_UUID",
		"./moeco.db",
	)

	err, errChan := MoecoSdk.Start(logger)
	if err != nil {
		logger.Errorf("%+v", err)
		panic(err)
	}

	for {
		err := <-errChan
		logger.Errorf("%+v", err)
	}
}
