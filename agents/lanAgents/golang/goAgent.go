/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/andy-zhangtao/humCICD/worker"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var giturl string
var branch string
var name string
var producer *nsq.Producer

const (
	BasePath = "/go/src/"
)

func nsqInit() {
	var errNum int
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(model.GoAgent)
	for {
		producer, _ = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
		err = producer.Ping()
		if err != nil {
			log.Output(model.GoAgent, "", logrus.Fields{"Ping Nsq Error": err}, logrus.ErrorLevel).Report()
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

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": producer.String()}).Info(model.GoAgent)
}

func valid() {
	if giturl == "" || branch == "" || name == "" {
		logrus.Error("git value or name value or branch value empty")
		os.Exit(-1)
	}
}

/*goAgent 构建Golang工程*/
func main() {
	app := cli.NewApp()
	app.Name = "goAgent"
	app.Usage = "clone & build golang project"
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
		cli.StringFlag{
			Name:        "name, n",
			Usage:       "The Hicd ID",
			Destination: &name,
		},
	}

	app.Action = buildAction
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}

}

func buildAction(c *cli.Context) error {
	defer log.Output(model.GoAgent, name, logrus.Fields{"msg": model.DefualtFinishFlag, "name": name}, logrus.InfoLevel).Report()
	nsqInit()
	valid()

	path, err := cloneGit(giturl, utils.ParsePath(giturl), branch)
	if err != nil {
		log.Output(model.GoAgent, name, logrus.Fields{"msg": fmt.Sprintf("Clone Error %s", err.Error()), "name": name}, logrus.ErrorLevel).Report()
		return err
	}

	configure, err := utils.GetConfigure(name)
	if err != nil {
		log.Output(model.GoAgent, name, logrus.Fields{"msg": fmt.Sprintf("Get Configure Error %s", err.Error()), "name": name}, logrus.ErrorLevel).Report()
		return err
	}

	ci := worker.CIWorker{
		Name:    name,
		WorkDir: path,
		Hicd:    configure.Configrue,
	}

	ci.Do()

	return nil
}

func cloneGit(url, name, branch string) (path string, err error) {

	// ref 需要提取project name
	if strings.HasPrefix(branch, "refs") {
		branch = strings.Split(branch, "refs/heads/")[1]
	}
	ref := "refs/remotes/origin/" + branch

	log.Output(model.GoAgent, name, logrus.Fields{"ref": ref, "path": name}, logrus.InfoLevel).Report()
	path = os.Getenv("GOPATH") + "/src/" + name
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
	})

	if err != nil {
		return
	}

	return
}