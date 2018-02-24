/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const _API_ = "/v1"
const ModuleName = "HUMCICD"

func main() {
	switch strings.ToLower(os.Getenv("HUM_DEBUG")) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.WithFields(logrus.Fields{"VERSION": getVersion()}).Info(ModuleName)

	r := mux.NewRouter()
	r.HandleFunc("/_ping", ping).Methods(http.MethodGet)
	r.HandleFunc(getAPIPath("/trigger"), trigger).Methods(http.MethodPost)
	logrus.Println(http.ListenAndServeTLS(":443", "server.crt", "server.key", r))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getVersion()))
}

func getVersion() string {
	return "v0.1"
}

func getAPIPath(path string) string {
	return _API_ + path
}

func trigger(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Read Body Error": err}).Error(ModuleName)
		return
	}

	logrus.WithFields(logrus.Fields{"Body": string(data)}).Info(ModuleName)
}
