/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/nsqio/go-nsq"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var giturl string
var branch string
var name string
var producer *nsq.Producer

func nsqInit() {
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(model.BuildAgent)
	producer, err = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{"Connect Nsq Error": err,}).Error(model.BuildAgent)
		os.Exit(-1)
	}

	err = producer.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Ping Nsq Error": err,}).Error(model.BuildAgent)
		os.Exit(-1)
	}

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": producer.String()}).Info(model.BuildAgent)
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
	nsqInit()
	valid()

	path, err := cloneGit(giturl, utils.ParsePath(giturl), branch)
	if err != nil {
		return err
	}

	defer logrus.WithFields(logrus.Fields{path: "Handler End"}).Info(model.BuildAgent)

	/*执行build*/
	if out, err := buildProject(path); err != nil {

		data, err := json.Marshal(model.OutEventMsg{
			Name:   name,
			Result: model.BuildFaild,
			Out:    fmt.Sprintf("ERR:[%s]\nLog:\n%s ", err.Error(), string(out)),
		})
		if err != nil {
			return err
		}
		producer.Publish(model.HicdOutTopic, data)
		return err
	} else {
		data, err := json.Marshal(model.OutEventMsg{
			Name:   name,
			Result: model.BuildSuc,
			Out:    fmt.Sprintf("Log:%s ", string(out)),
		})
		if err != nil {
			return err
		}
		producer.Publish(model.HicdOutTopic, data)
	}

	return nil
}

func cloneGit(url, name, branch string) (path string, err error) {
	logrus.WithFields(logrus.Fields{"ref": plumbing.ReferenceName("refs/heads/" + branch), "path": name}).Info(model.BuildAgent)
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

func buildProject(path string) (result []byte, err error) {
	var out, stderr bytes.Buffer
	var cmd *exec.Cmd
	err = os.Chdir(path)
	if err != nil {
		result = []byte(fmt.Sprintf("%s \n %s", out.String(), err.Error()))
		return
	}

	_, err = os.Stat(path + "/Makefile")
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Println("Not Has Makefile")
			cmd = exec.Command("go", []string{"build", "-v"}...)
		}
	} else {
		logrus.Println("Has Makefile")
		cmd = exec.Command("make", "all")
	}

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	//
	//out, err := cmd.CombinedOutput()
	err = cmd.Run()
	if err != nil {
		result = []byte(fmt.Sprintf("%s\n%s\n%s", out.String(), stderr.String(), err.Error()))
		return
	}

	result = []byte(fmt.Sprintf("Output:\n%s \nErr:\n%s \n", out.String(), stderr.String()))

	return
}
