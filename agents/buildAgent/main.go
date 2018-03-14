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

	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/fsouza/go-dockerclient"
	"github.com/nsqio/go-nsq"
)

var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message

/*buildAgent 从NSQ读取工程解析后的数据，然后执行构建任务*/

type BuildAgent struct {
	Name        string
	NsqEndpoint string
	Client      *docker.Client
}

type gitInfo struct {
	Git    string
	Branch string
}

func (this *BuildAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *BuildAgent) Run() {
	if err := this.checkRun(); err != nil {
		logrus.WithFields(logrus.Fields{"BuildAgent CheckRun Failed": err}).Error(this.Name)
		return
	}

	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.GitAgentTopic, this.Name, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name}).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(logrus.Fields{"WorkChan": "Listen..."}).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
			msg := model.GitConfigure{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			go this.handleBuild(msg)

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&BuildAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint}).Info(this.Name)
	<-r.StopChan
}

// checkRun 检查是否具备运行环境
// 包括检查是否具备docker运行条件
func (this *BuildAgent) checkRun() error {
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
	}

	summry, err := this.Client.ListImages(docker.ListImagesOptions{
		All:    false,
		Filter: fmt.Sprintf("reference=%s", model.GoImage),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"List Image Error": err}).Error(this.Name)
		return err
	}

	if len(summry) == 0 {
		logrus.WithFields(logrus.Fields{"Is Has goAgent": false, "Pull Image": "..."}).Info(this.Name)
		this.Client.PullImage(docker.PullImageOptions{
			Context:    context.Background(),
			Repository: model.GoImage,
			Tag:        "latest",
		}, docker.AuthConfiguration{})
		//this.Client.ImagePull(context.Background(), model.GoImage, types.ImagePullOptions{})
	} else {
		logrus.WithFields(logrus.Fields{"Is Has goAgent": true}).Info(this.Name)
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

func (this *BuildAgent) handleBuild(msg model.GitConfigure) {
	logrus.WithFields(logrus.Fields{"Name": msg.Name, "Configrue": msg.Configrue}).Info(this.Name)
	switch msg.Configrue.Kind {
	case "golang":
		this.buildGolang(msg)
	}
}

func (this *BuildAgent) buildGolang(msg model.GitConfigure) {
	/*1. 构建golang容器*/
	opt := model.BuildOpts{
		Client: this.Client,
		DockerOpt: []model.DockerOpts{model.DockerOpts{
			Img: "vikings/goagent",
			Cmd: fmt.Sprintf("-g %s -b %s -n %s", msg.GitUrl, msg.Branch, msg.Name),
		}},
	}
	err := utils.CreateContainer(opt)
	if err != nil {
		logrus.Error(err)
	}
}

func main() {
	bagent := BuildAgent{
		Name:        model.BuildAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	bagent.Run()
}
