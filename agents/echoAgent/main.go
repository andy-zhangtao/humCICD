/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andy-zhangtao/gogather/tools"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message
// projectMsg 保存每个工程所对应的日志
var projectMsg map[string]string

type EchoAgent struct {
	Name        string
	NsqEndpoint string
}

func (this *EchoAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *EchoAgent) Run() {

	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.HicdOutTopic, this.Name, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name}).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(logrus.Fields{"WorkChan": "Listen..."}).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
			msg := model.OutEventMsg{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			go this.handlerOutput(msg)

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&EchoAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint}).Info(this.Name)
	<-r.StopChan
}

func (this *EchoAgent) handlerOutput(msg model.OutEventMsg) {
	logrus.WithFields(logrus.Fields{"Name": msg.Name, "Project": msg.Project, "Result": msg.Result}).Info(this.Name)
	msg.Out = strings.Replace(msg.Out, "\n", "<br/>", -1)
	logrus.Print(msg.Out)

	if projectMsg[msg.Project] != "" {
		projectMsg[msg.Project] += msg.Out + "<br/>"
	} else {
		projectMsg[msg.Project] = msg.Out
	}

	switch msg.Result {
	case model.BuildSuc:
		break
	case model.BuildFaild:
		e := tools.Email{
			Host:     os.Getenv(model.EnvEmailHost),
			Username: os.Getenv(model.EnvEmailUser),
			Password: os.Getenv(model.EnvEmailPass),
			Port:     587,
			Dest:     []string{os.Getenv(model.EnvEmailDest)},
			Content:  projectMsg[msg.Project],
			Header:   fmt.Sprintf("HICD [%s] Report", msg.Project),
		}
		if err := e.SendEmail(); err != nil {
			logrus.WithFields(logrus.Fields{"Send Email Error": err}).Error(this.Name)
		}
	}
}

/*EchoAgent 从NSQ读取所有成功或者失败的信息*/
func main() {
	eagent := EchoAgent{
		Name:        model.EchoAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	projectMsg = make(map[string]string)
	eagent.Run()
}
