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

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
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
var producer *nsq.Producer

func nsqInit() {
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		log.Output(model.GitAgent, "", logrus.Fields{"Env Empty": model.EnvNsqdEndpoint}, logrus.ErrorLevel).Report()
		// logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	log.Output(model.GitAgent, "", logrus.Fields{"Connect NSQ": nsq_endpoint}, logrus.DebugLevel)
	producer, err = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		log.Output(model.GitAgent, "", logrus.Fields{"Connect Nsq Error": err}, logrus.ErrorLevel).Report()
		os.Exit(-1)
	}

	err = producer.Ping()
	if err != nil {
		log.Output(model.GitAgent, "", logrus.Fields{"Ping Nsq Error": err}, logrus.ErrorLevel).Report()
		os.Exit(-1)
	}

	log.Output(model.GitAgent, "", logrus.Fields{"Connect Nsq Succes": producer.String()}, logrus.InfoLevel)
}

func valid() {
	if giturl == "" || branch == "" {
		log.Output(model.GitAgent, branch, logrus.Fields{"Parameter Error": "git value or branch value empty"}, logrus.ErrorLevel)
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
			Name:        "git, g",
			Usage:       "The Git URL",
			Destination: &giturl,
		},
		cli.StringFlag{
			Name:        "branch, b",
			Usage:       "The Git Branch Name",
			Destination: &branch,
		},
		// cli.StringFlag{
		// 	Name:        "name, n",
		// 	Usage:       "Hicd ID",
		// 	Destination: &name,
		// },
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
	configrue, err := cloneGit(giturl, parseName(giturl), branch)
	if err != nil {
		return err
	}

	return sendConfigure(configrue)
}

func cloneGit(url, name, branch string) (configure *model.HICD, err error) {
	ref := ""
	if strings.HasPrefix(branch, "refs") {
		ref = branch
	} else {
		ref = "refs/remotes/origin/" + branch
	}
	log.Output(model.GitAgent, name, logrus.Fields{"ref": plumbing.ReferenceName(ref), "msg": ref}, logrus.InfoLevel).Report()

	_, err = git.PlainClone("/tmp/"+name, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(ref),
	})

	if err != nil {
		return
	}

	configure, err = parseConfigure("/tmp/" + name)
	if err != nil {
		return
	}

	log.Output(model.GitAgent, name, logrus.Fields{"msg": fmt.Sprintf("language:[%s]", configure.Language)}, logrus.InfoLevel).Report()
	return
}

// parseName 从Git URL中提取工程名
// 例如从https://github.com/andy-zhangtao/humCICD.git中提取出humCICD
func parseName(url string) (string) {
	gitName := strings.Split(url, "/")
	project = strings.Split(gitName[len(gitName)-1], ".")[0]
	log.Output(model.GitAgent, project, logrus.Fields{"Process": fmt.Sprintf("GitAgent Will Clone [%s]\n", project)}, logrus.InfoLevel)
	return project
}

// parseConfigrue 解析工程中的.hicd文件
// path 工程路径
func parseConfigure(path string) (configure *model.HICD, err error) {
	configure = new(model.HICD)
	fileName := path + "/.hicd.toml"
	_, err = os.Open(fileName)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("Open .hicd.toml error[%s]", err))
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Read .hicd.toml error[%s]", err))
	}

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
		Configrue: *configure,
	}

	data, err := json.Marshal(&hc)
	if err != nil {
		return err
	}

	return producer.Publish(model.GitAgentTopic, data)
}
