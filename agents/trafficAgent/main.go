/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"


"github.com/andy-zhangtao/humCICD/model"
"github.com/andy-zhangtao/humCICD/utils"
"github.com/fsouza/go-dockerclient"
"github.com/nsqio/go-nsq"
"github.com/sirupsen/logrus"

)

var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message

type TrafficAgent struct {
	Name        string
	NsqEndpoint string
	Client      *docker.Client
}

func (this *TrafficAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *TrafficAgent) Run() {

	if err := this.checkRun(); err != nil {
		logrus.WithFields(logrus.Fields{"TrafficAgent CheckRun Failed": err}).Error(this.Name)
		return
	}

	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.TAGQUEUE, this.Name, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name}).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(logrus.Fields{"WorkChan": "Listen..."}).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
			msg := model.TagEventMsg{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			go this.handlerGit(msg)

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&TrafficAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint}).Info(this.Name)
	<-r.StopChan
}

// handlerGit 处理GitHub发来的通知消息
func (this *TrafficAgent) handlerGit(msg model.TagEventMsg) {
	logrus.WithFields(logrus.Fields{"Create gitAgent": fmt.Sprintf("-g %s -b %s -n %s", msg.GitURL, msg.Branch, msg.Name)}).Info(this.Name)
	opt := model.BuildOpts{
		Client: this.Client,
		DockerOpt: []model.DockerOpts{model.DockerOpts{
			Img:  "vikings/gitagent:latest",
			Cmd:  fmt.Sprintf("-g %s -b %s -n %s", msg.GitURL, msg.Branch, msg.Name),
			Name: "gitagent-" + msg.Name,
			Env:  map[string]string{model.EnvNsqdEndpoint: os.Getenv(model.EnvNsqdEndpoint)},
		}},
	}
	err := utils.CreateContainer(opt)
	if err != nil {
		logrus.Error(err)
	}

}

// checkRun 检查是否具备运行环境
// 包括检查是否具备docker运行条件
func (this *TrafficAgent) checkRun() error {
	/*check docker runtime*/
	if cli, err := checkDocker(); err != nil {
		return errors.New(fmt.Sprintf("Check Docker Error [%v]", err))
	} else {
		this.Client = cli
		logrus.WithFields(logrus.Fields{"Docker Version": this.Client.Version()}).Info(this.Name)
		//filter := filters.NewArgs()
		//filter.Add("reference", model.GitImage)
		//summry, err := this.Client.ImageList(context.Background(), types.ImageListOptions{
		//	All:     false,
		//	Filters: filter,
		//})
		//if err != nil {
		//	logrus.WithFields(logrus.Fields{"List Image Error": err}).Error(this.Name)
		//	return err
		//}
		//
		//if len(summry) == 0 {
		//	logrus.WithFields(logrus.Fields{"Is Has gitAgent": false, "Pull Image": "..."}).Info(this.Name)
		//	this.Client.ImagePull(context.Background(), model.GitImage, types.ImagePullOptions{})
		//} else {
		//	logrus.WithFields(logrus.Fields{"Is Has gitAgent": true}).Info(this.Name)
		//}
	}

	return nil
}

func checkDocker() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	//_, err = client.Ping(context.Background())
	return
}

func main() {

	tagent := TrafficAgent{
		Name:        model.TrafficAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	tagent.Run()
}
