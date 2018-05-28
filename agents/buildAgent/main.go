/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"context"
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
	logrus.WithFields(log.Z().Fields(logrus.Fields{"HandleMessage": string(m.Body)})).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *BuildAgent) Run() {
	if err := this.checkRun(); err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"BuildAgent CheckRun Failed": err})).Error(this.Name)
		return
	}

	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.GitConfIDTopic, this.Name, cfg)
	if err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name})).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"WorkChan": "Listen..."})).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(log.Z().Fields(logrus.Fields{"BuildMsg": string(m.Body)})).Info(this.Name)
			// msg := model.GitConfigure{}
			//
			// err = json.Unmarshal(m.Body, &msg)
			// if err != nil {
			// 	logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
			// 	continue
			// }

			go this.handleBuild(string(m.Body))

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&BuildAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint})).Info(this.Name)
	<-r.StopChan
}

// checkRun 检查是否具备运行环境
// 包括检查是否具备docker运行条件
func (this *BuildAgent) checkRun() error {
	/*check docker runtime*/
	if cli, err := checkDocker(); err != nil {
		return log.Z().Error(fmt.Sprintf("Check Docker Error [%v]", err))
	} else {
		this.Client = cli
		env, err := this.Client.Version()
		if err != nil {
			return err
		}
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Docker Version": env.Get("Version")})).Info(this.Name)
	}

	summry, err := this.Client.ListImages(docker.ListImagesOptions{
		Filters: map[string][]string{
			"reference": {model.GoImage},
		},
		All: false,
		// Filter: fmt.Sprintf("reference=%s", model.GoImage),
	})
	if err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"List Image Error": err})).Error(this.Name)
		return err
	}

	if len(summry) == 0 {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Is Has goAgent": false, "Pull Image": "..."})).Info(this.Name)
		this.Client.PullImage(docker.PullImageOptions{
			Context:    context.Background(),
			Repository: model.GoImage,
			Tag:        "latest",
		}, docker.AuthConfiguration{})
		// this.Client.ImagePull(context.Background(), model.GoImage, types.ImagePullOptions{})
	} else {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Is Has goAgent": true})).Info(this.Name)
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

func (this *BuildAgent) handleBuild(msgid string) {
	configure, err := utils.GetConfigure(msgid)
	if err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Get Configrue Err": err})).Info(this.Name)
		return
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{"Name": configure.Name, "Configrue": configure})).Info(this.Name)
	switch configure.Configrue.Language {
	case "golang":
		this.buildGolang(configure)
	}
}

func (this *BuildAgent) buildGolang(msg *model.GitConfigure) {
	/*1. 构建容器*/

	envMap := map[string]string{
		model.EnvNsqdEndpoint: os.Getenv(model.EnvNsqdEndpoint),
		model.EnvDataAgent:    os.Getenv(model.EnvDataAgent),
	}

	if !msg.Configrue.Env.Skip {
		for _, env := range msg.Configrue.Env.Var {
			for key, value := range env {
				if key != "" && value != "" {
					envMap[key] = value
				}
			}
		}
	}

	cmd := fmt.Sprintf("-t %s -g %s -b %s -n %s", log.Z().MyTrack(), msg.GitUrl, msg.Branch, msg.Name)
	logrus.WithFields(log.Z().Fields(logrus.Fields{"create contianer": cmd})).Info(this.Name)
	opt := model.BuildOpts{
		Client: this.Client,
		DockerOpt: []model.DockerOpts{model.DockerOpts{
			Img: "vikings/goagent",
			Cmd: cmd,
			Env: envMap,
		}},
	}

	if msg.Configrue.After.Usedocker {
		opt.DockerOpt[0].Binds = []string{"/var/run/docker.sock:/var/run/docker.sock"}
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
