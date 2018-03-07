/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/fsouza/go-dockerclient"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/24.
var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message

type BuildAgent struct {
	Name        string
	NsqEndpoint string
	Client      *docker.Client
}

type BuildOpts struct {
	DockerOpt []DockerOpts `json:"docker_opt"`
}

type DockerOpts struct {
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Port []int             `json:"port"`
	Env  map[string]string `json:"env"`
	Cmd  string            `json:"cmd"`
}

// HCIDOpt CI/CD的配置数据
type HCIDOpt struct {
}

func (this *BuildAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	//workerHome[this.Name] <- m
	workerChan <- m
	return nil
}

func (this *BuildAgent) Run() {
	if err := this.checkRun(); err != nil {
		logrus.WithFields(logrus.Fields{"BuildAgent CheckRun Failed": err}).Error(this.Name)
		return
	}

	workerChan = make(chan *nsq.Message)

	//workerHome = make(map[string]chan *nsq.Message)
	//workerHome[this.Name] = workerChan

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

			logrus.WithFields(logrus.Fields{"Kind": msg.Kind, "Branch": msg.Branch, "GitURL": msg.GitURL, "Tag": msg.Tag}).Info(this.Name)

			go this.handleBuild()

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

// handleBuild 处理构建请求
// 创建docker容器
// 从github获取指定版本的代码进行构建
func (this *BuildAgent) handleBuild() {
	/*解析HCID配置数据*/
	//parseConfigrue
}

// checkRun 检查是否具备运行环境
// 包括检查是否具备docker运行条件
func (this *BuildAgent) checkRun() error {
	/*check docker runtime*/
	if cli, err := checkDocker(); err != nil {
		return errors.New(fmt.Sprintf("Check Docker Error [%v]", err))
	} else {
		this.Client = cli
		logrus.WithFields(logrus.Fields{"Docker Check": true}).Info(this.Name)
	}

	return nil
}

func checkDocker() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		return
	}

	err = client.Ping()
	return
}

// buildContainer 按照配置文件中的约束关系进行容器创建
func (this *BuildAgent) buildContainer(do DockerOpts) error {
	pb := make(map[docker.Port][]docker.PortBinding)
	if len(do.Port) > 0 {
		for _, p := range do.Port {
			pb[docker.Port(fmt.Sprintf("%d/tcp", p))] = []docker.PortBinding{
				docker.PortBinding{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", p)},
			}
		}
	}

	logrus.WithFields(logrus.Fields{"Port": pb}).Info(this.Name)

	container, err := this.Client.CreateContainer(docker.CreateContainerOptions{
		Name: do.Name,
		Config: &docker.Config{
			Image:        do.Img,
			AttachStdout: true,
			AttachStdin:  true,
		},
		HostConfig: &docker.HostConfig{
			PortBindings: pb,
		},
		NetworkingConfig: &docker.NetworkingConfig{},
		Context:          context.Background(),
	})

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Name": do.Name, "ID": container.ID}).Info(this.Name)

	return nil
}

// parseConfigrue 解析HCID配置文件
// 1. clone 指定的git地址
// 2. checkout 指定分支
// 3. 判断是否存在.hcid.yml文件
// 4. 解析.hcid.yml文件
func (this *BuildAgent) parseConfigrue() (opt *HCIDOpt, err error) {
	return nil, nil
}
