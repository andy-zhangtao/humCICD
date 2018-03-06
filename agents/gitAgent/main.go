/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
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
	return
}

func parseName(url string) (name string) {
	gitName := strings.Split(url, "/")
	name = strings.Split(gitName[len(gitName)-1], ".")[0]
	fmt.Printf("GitAgent Will Clone [%s]\n", name)
	return
}
