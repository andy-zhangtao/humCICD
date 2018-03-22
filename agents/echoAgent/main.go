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
	"time"

	"github.com/andy-zhangtao/gogather/tools"
	"github.com/andy-zhangtao/humCICD/influx"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message

type EchoAgent struct {
	Name        string
	NsqEndpoint string
}

func (this *EchoAgent) HandleMessage(m *nsq.Message) error {
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
			// logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
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
	logrus.WithFields(logrus.Fields{"Name": msg.Name, "Project": msg.Project, "Result": msg.Result}).Info(model.EchoAgent)

	if msg.Project == "" {
		// 如果Project为空,则是和业务无关联的日志. 不需要发送这些日志
		return
	}

	if msg.Out == "" {
		return
	}

	msg.Out = strings.Replace(msg.Out, "\n", "<br/>", -1)

	err := influx.Insert(msg.Project, logrus.Fields{"name": msg.Project}, logrus.Fields{"log": msg.Out})
	if err != nil {
		logrus.WithFields(logrus.Fields{"Save InfluxDB Error": err, "name": msg.Project}).Error(model.EchoAgent)
		return
	}

	if msg.Out == model.DefualtFinishFlag {
		// 任务结束,需要发送邮件
		logrus.WithFields(logrus.Fields{"Query InfluxDB": msg.Project, "End": true}).Info(model.EchoAgent)
		err := sendEmail(msg.Project)
		if err != nil {
			return
		}

		err = influx.Destory(msg.Project)
		if err != nil {
			return
		}
	}
}

func sendEmail(project string) error {

	runLog, err := influx.Query(project)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Query Log Error": err}).Error(model.EchoAgent)
		return err
	}

	content := ""
	for _, l := range runLog {
		content += fmt.Sprintf(" [%d] %s <br/>", time.Unix(l.Timestamp, 0), l.Message)
	}

	logrus.WithFields(logrus.Fields{"Email": content}).Info(model.EchoAgent)

	e := tools.Email{
		Host:     os.Getenv(model.EnvEmailHost),
		Username: os.Getenv(model.EnvEmailUser),
		Password: os.Getenv(model.EnvEmailPass),
		Port:     587,
		Dest:     []string{os.Getenv(model.EnvEmailDest)},
		Content:  content,
		Header:   fmt.Sprintf("HICD [%s] Report", project),
	}
	if err := e.SendEmail(); err != nil {
		logrus.WithFields(logrus.Fields{"Send Email Error": err}).Error(model.EchoAgent)
		return err
	}

	return nil
}

/*EchoAgent 从NSQ读取所有成功或者失败的信息*/
func main() {
	eagent := EchoAgent{
		Name:        model.EchoAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	eagent.Run()
}
