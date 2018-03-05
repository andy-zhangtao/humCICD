/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

/*
gitAgent 从Github上面拉取指定工程
然后解析工程中HICD的配置数据
*/
func main() {
	app := cli.NewApp()
	app.Name = "gitAgent"
	app.Usage = "clone & parse HICD configure"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
