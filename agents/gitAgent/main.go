/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/nsqio/go-nsq"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

/*
gitAgent 从Github上面拉取指定工程
然后解析工程中HICD的配置数据
*/
var giturl string
// project 工程名称
var project string
var branch string
var email string
var producer *nsq.Producer
var track string

func nsqInit() {
	var errNum int
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		log.Output(model.GitAgent, "", logrus.Fields{"Env Empty": model.EnvNsqdEndpoint}, logrus.ErrorLevel).Report()
		os.Exit(-1)
	}

	log.Output(model.GitAgent, "", logrus.Fields{"Connect NSQ": nsq_endpoint}, logrus.DebugLevel)
	for {
		producer, _ = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
		err = producer.Ping()
		if err != nil {
			log.Output(model.GitAgent, "", logrus.Fields{"Ping Nsq Error": err}, logrus.ErrorLevel).Report()
			errNum ++
		}

		if err == nil {
			break
		}

		if errNum >= 20 {
			os.Exit(-1)
		}
		time.Sleep(time.Second * 5)
	}

	log.Output(model.GitAgent, "", logrus.Fields{"Connect Nsq Succes": producer.String()}, logrus.InfoLevel)
}

func valid() {
	if giturl == "" || branch == "" || track == "" {
		log.Output(model.GitAgent, branch, logrus.Fields{"Parameter Error": "git value or branch value or track value empty"}, logrus.ErrorLevel)
		os.Exit(-1)
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "gitAgent"
	app.Usage = "clone & parse HICD configure"
	app.Version = "v0.1.0"
	app.Author = "andy zhang"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "track, t",
			Usage:       "The Log Track ID",
			Destination: &track,
		},
		cli.StringFlag{
			Name:        "git, g",
			Usage:       "The Git URL",
			Destination: &giturl,
		},
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "The Git Branch Name",
			Destination: &branch,
		},
		cli.StringFlag{
			Name:        "email, e",
			Usage:       "Pusher Email",
			Destination: &email,
		},
	}

	app.Action = parseAction
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func parseAction(c *cli.Context) error {
	nsqInit()
	valid()
	log.Z().AddID(track)
	configrue, err := cloneGit(giturl, parseName(giturl), branch)
	if err != nil {
		return err
	}

	return sendConfigure(configrue)
}

func cloneGit(url, name, branch string) (configure *model.HICD, err error) {
	project = utils.ParsePath(url)
	ref := ""
	if strings.HasPrefix(branch, "refs") {
		ref = branch
	} else {
		ref = "refs/remotes/origin/" + branch
	}

	log.Output(model.GitAgent, project, logrus.Fields{"ref": plumbing.ReferenceName(ref), "msg": ref}, logrus.InfoLevel).Report()

	logrus.WithFields(log.Z().Fields(logrus.Fields{"ref": plumbing.ReferenceName(ref), "msg": ref})).Info(model.GitAgent)
	_, err = git.PlainClone("/tmp/"+name, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(ref),
	})

	if err != nil {
		logrus.Error(log.Z().Error(err.Error()))
		return
	}

	configure, err = parseConfigure("/tmp/" + name)
	if err != nil {
		logrus.Error(log.Z().Error(err.Error()))
		log.Output(model.GitAgent, project, logrus.Fields{"msg": err.Error()}, logrus.ErrorLevel).Report()
		// 如果读取.hicd.toml失败,此时应该退出
		log.Output(model.GitAgent, project, logrus.Fields{"msg": model.DefualtFinishFlag}, logrus.ErrorLevel).Report()
		return
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{"msg": fmt.Sprintf("language:[%s]", configure.Language)})).Info(model.GitAgent)
	log.Output(model.GitAgent, project, logrus.Fields{"msg": fmt.Sprintf("language:[%s]", configure.Language)}, logrus.InfoLevel).Report()
	return
}

// parseName 从Git URL中提取工程名
// 例如从https://github.com/andy-zhangtao/humCICD.git中提取出humCICD
func parseName(url string) (string) {
	gitName := strings.Split(url, "/")
	name := strings.Split(gitName[len(gitName)-1], ".")[0]
	logrus.WithFields(log.Z().Fields(logrus.Fields{"Process": fmt.Sprintf("GitAgent Will Clone [%s]\n", name)})).Info(model.GitAgent)
	log.Output(model.GitAgent, name, logrus.Fields{"Process": fmt.Sprintf("GitAgent Will Clone [%s]\n", name)}, logrus.InfoLevel)
	return name
}

// parseConfigrue 解析工程中的.hicd文件
// path 工程路径
func parseConfigure(path string) (configure *model.HICD, err error) {
	configure = new(model.HICD)
	fileName := path + "/.hicd.toml"
	_, err = os.Open(fileName)
	if os.IsNotExist(err) {
		logrus.Error(log.Z().Error(fmt.Sprintf("Open .hicd.toml error[%s]", err)))
		return nil, errors.New(fmt.Sprintf("Open .hicd.toml error[%s]", err))
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Error(log.Z().Error(fmt.Sprintf("Read .hicd.toml error[%s]", err)))
		return nil, errors.New(fmt.Sprintf("Read .hicd.toml error[%s]", err))
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{".hcid": string(data)})).Info(model.GitAgent)
	log.Output(model.GitAgent, project, logrus.Fields{".hcid": string(data)}, logrus.InfoLevel)

	err = toml.Unmarshal(data, configure)
	if err != nil {
		return nil, err
	}
	return
}

// sendConfigure 发送配置消息
func sendConfigure(configure *model.HICD) error {

	hc := model.GitConfigure{
		Name:      project,
		GitUrl:    giturl,
		Branch:    branch,
		Email:     email,
		Configrue: *configure,
	}

	data, err := json.Marshal(&hc)
	if err != nil {
		return err
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{"configure": hc})).Info(model.GitAgent)
	log.Output(model.GitAgent, model.DefualtEmptyProject, logrus.Fields{"configure": hc}, logrus.InfoLevel)
	return producer.Publish(model.GitAgentTopic, data)
}
