/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/andy-zhangtao/humCICD/log"
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
			// msg := model.TagEventMsg{}
			msg := model.EventMsg{}
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
func (this *TrafficAgent) handlerGit(msg model.EventMsg) {
	opt := model.BuildOpts{
		Client: this.Client,
	}
	switch msg.Kind {
	case model.PushEventType:
		m := msg.Msg.(map[string]interface{})
		gitURL := m["git_url"].(string)
		branch := m["branch"].(string)
		name := m["name"].(string)
		email := m["email"].(string)
		log.Output(this.Name, branch, logrus.Fields{"Create gitAgent": fmt.Sprintf("-g %s -b %s", gitURL, branch)}, logrus.InfoLevel)
		opt.DockerOpt = []model.DockerOpts{model.DockerOpts{
			Img: "vikings/gitagent:latest",
			// Cmd:  fmt.Sprintf("-g %s -b %s -n %s", gitURL, branch, name),
			Cmd:  fmt.Sprintf("-g %s -b %s -e %s", gitURL, branch, email),
			Name: "gitagent-" + name,
			Env:  map[string]string{model.EnvNsqdEndpoint: os.Getenv(model.EnvNsqdEndpoint)},
		}}
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
		env, err := this.Client.Version()
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{"Docker Version": env.Get("Version")}).Info(this.Name)
		summry, err := this.Client.ListImages(docker.ListImagesOptions{
			Filters: map[string][]string{
				"reference": {model.GitImage},
			},
			All: false,
			// Filter: fmt.Sprintf("reference=%s", model.GitImage),
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{"List Image Error": err}).Error(this.Name)
			return err
		}

		if len(summry) == 0 {
			logrus.WithFields(logrus.Fields{"Is Has gitAgent": false, "Pull Image": "..."}).Info(this.Name)
			this.Client.PullImage(docker.PullImageOptions{
				Context:    context.Background(),
				Repository: model.GitImage,
				Tag:        "latest",
			}, docker.AuthConfiguration{})
			// this.Client.ImagePull(context.Background(), model.GoImage, types.ImagePullOptions{})
		} else {
			logrus.WithFields(logrus.Fields{"Is Has gitAgent": true}).Info(this.Name)
		}
	}

	return nil
}

func checkDocker() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	err = client.Ping()
	return
}

func main() {

	tagent := TrafficAgent{
		Name:        model.TrafficAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	tagent.Run()
}
