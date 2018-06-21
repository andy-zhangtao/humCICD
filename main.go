/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/andy-zhangtao/_hulk_client"
	"github.com/andy-zhangtao/humCICD/git"
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/gorilla/mux"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

const _API_ = "/v1"
const ModuleName = "HUMCICD"

const ServiceName = "HICD_MAIN_AGENT"
const ServiceVersion = "v1.0.0"
const ServiceResume = "HICD_MAIN_AGENT监听来自于Github的Pull请求,并会进行初步解析处理。然后将此数据放入NSQ中"

type NsqBridge struct {
	producer *nsq.Producer
}

var nb *NsqBridge

func main() {
	defer func() {
		_hulk_client.UnRegister(ServiceName, ServiceVersion)
	}()

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
	r.HandleFunc(getAPIPath("/tag/trigger"), trigger).Methods(http.MethodPost)
	r.HandleFunc(getAPIPath("/push/trigger"), push).Methods(http.MethodPost)
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

func push(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Output(ModuleName, model.DefualtEmptyProject, logrus.Fields{"Read Body Error": err}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(ModuleName, model.DefualtEmptyProject, logrus.Fields{"Body": string(data)}, logrus.DebugLevel)

	pushEvent := git.GitHubPush{}

	err = json.Unmarshal(data, &pushEvent)
	if err != nil {
		log.Output(ModuleName, model.DefualtEmptyProject, logrus.Fields{"Unmarshal Body Error": err}, logrus.ErrorLevel).Report()
		return
	}

	push := model.PushEventMsg{
		GitURL: pushEvent.Repository.Clone_url,
		Branch: pushEvent.Ref,
		Name:   pushEvent.Repository.Full_name,
		Email:  pushEvent.Pusher.Email,
	}

	msg := model.EventMsg{
		Kind:  model.PushEventType,
		Email: pushEvent.Pusher.Email,
		Msg:   push,
	}

	m, err := json.Marshal(msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Marshal Body Error": err}).Error(ModuleName)
		return
	}

	makeMsg(model.TAGQUEUE, string(m))
}

func trigger(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Read Body Error": err}).Error(ModuleName)
		return
	}

	logrus.WithFields(logrus.Fields{"Body": string(data)}).Debug(ModuleName)

	tagEvent := git.TagEvent{}

	err = json.Unmarshal(data, &tagEvent)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Unmarshal Body Error": err}).Error(ModuleName)
		return
	}

	branch := strings.Split(tagEvent.Ref, "-")
	if len(branch) == 1 {
		branch = append(branch, "master")
	}

	msg := model.TagEventMsg{
		Kind:   "tag",
		GitURL: tagEvent.Repository.Clone_url,
		Tag:    branch[0],
		Branch: branch[1],
		Name:   tagEvent.Repository.Name,
	}

	m, err := json.Marshal(msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Marshal Body Error": err}).Error(ModuleName)
		return
	}

	makeMsg(model.TAGQUEUE, string(m))

}

func init() {
	_hulk_client.Register(ServiceName, ServiceVersion, ServiceResume)
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(ModuleName)

	producer, err := nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{"Connect Nsq Error": err,}).Panic(ModuleName)
	}

	nb = &NsqBridge{
		producer: producer,
	}

	err = producer.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Ping Nsq Error": err,}).Panic(ModuleName)
	}

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": producer.String(),}).Info(ModuleName)

}

func makeMsg(topic, msg string) error {
	logrus.WithFields(logrus.Fields{"Topic": topic, "Msg": msg}).Info(ModuleName)
	return nb.producer.Publish(topic, []byte(msg))
}
