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
	"os/exec"

	"github.com/Sirupsen/logrus"
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

	/*执行build*/
	if err = buildProject(path); err != nil {

		data, err := json.Marshal(model.ErrEventMsg{
			Name: name,
			Err:  err.Error(),
		})
		if err != nil {
			return err
		}
		producer.Publish(model.HicdErrTopic, data)
		return err
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

func buildProject(path string) error {
	var cmd *exec.Cmd
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path + "/Makefile"); os.IsExist(err) {
		/*存在Makefile*/
		cmd = exec.Command("make")
	} else {
		/*不存在Makefile*/
		cmd = exec.Command("go", "build")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}

	derr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}

	logrus.Println(string(out))
	if len(derr) > 0 {
		logrus.Errorln(string(derr))
		return errors.New(string(derr))
	}
	return nil
}
