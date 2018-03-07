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
	"github.com/andy-zhangtao/humCICD/agents/models"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
)

/*
gitAgent 从Github上面拉取指定工程
然后解析工程中HICD的配置数据
*/
var giturl string

func main() {
	app := cli.NewApp()
	app.Name = "gitAgent"
	app.Usage = "clone & parse HICD configure"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "git, g",
			Usage:       "The Git URL",
			Destination: &giturl,
		},
	}

	//app.Action = func(c *cli.Context) error {
	//	fmt.Println("boom! I say!" + giturl)
	//	return nil
	//}

	app.Action = parseAction
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func parseAction(c *cli.Context) error {
	err := cloneGit(giturl, parseName(giturl))
	if err != nil {
		return err
	}

	return nil
}

func cloneGit(url, name string) (err error) {
	_, err = git.PlainClone("/tmp/"+name, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	hc, err := parseConfigure("/tmp/" + name)
	if err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{"configrue": hc}).Info("gitAgent")
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
func parseConfigure(path string) (configure *models.HicdConfigure, err error) {
	configure = new(models.HicdConfigure)
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
