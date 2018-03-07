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

	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/nsqio/go-nsq"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
)

/*
gitAgent 从Github上面拉取指定工程
然后解析工程中HICD的配置数据
*/
var giturl string
var name string
var producer *nsq.Producer

func nsqInit() {
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(model.GitAgent)
	producer, err = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{"Connect Nsq Error": err,}).Error(model.GitAgent)
		os.Exit(-1)
	}

	err = producer.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Ping Nsq Error": err,}).Error(model.GitAgent)
		os.Exit(-1)
	}

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": producer.String()}).Info(model.GitAgent)
}

func valid() {
	if giturl == "" || name == "" {
		logrus.Error("git value or name value empty")
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
			Name:        "name, n",
			Usage:       "Hicd ID",
			Destination: &name,
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
	configrue, err := cloneGit(giturl, parseName(giturl))
	if err != nil {
		return err
	}

	return sendConfigure(configrue)
}

func cloneGit(url, name string) (configure *model.HicdConfigure, err error) {
	_, err = git.PlainClone("/tmp/"+name, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	configure, err = parseConfigure("/tmp/" + name)
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{"configrue": configure}).Info("gitAgent")
	return
}

func parseName(url string) (name string) {
	gitName := strings.Split(url, "/")
	name = strings.Split(gitName[len(gitName)-1], ".")[0]
	fmt.Printf("GitAgent Will Clone [%s]\n", name)
	return
}

// parseConfigrue 解析工程中的.hicd文件
// path 工程路径
func parseConfigure(path string) (configure *model.HicdConfigure, err error) {
	configure = new(model.HicdConfigure)
	fileName := path + "/.hicd"
	_, err = os.Open(fileName)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("Open .hicd error[%s]", err))
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Read .hicd error[%s]", err))
	}

	logrus.WithFields(logrus.Fields{".hcid": string(data)}).Info("gitAgent")

	err = json.Unmarshal(data, configure)
	if err != nil {
		return nil, err
	}
	return
}

// sendConfigure 发送配置消息
func sendConfigure(configure *model.HicdConfigure) error {
	type HicdConfigure struct {
		Name      string              `json:"name"`
		Configrue model.HicdConfigure `json:"configrue"`
	}

	hc := HicdConfigure{
		Name:      name,
		Configrue: *configure,
	}

	data, err := json.Marshal(&hc)
	if err != nil {
		return err
	}

	return producer.Publish(model.GitAgentTopic, data)
}
