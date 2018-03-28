/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package worker

import (
	"fmt"
	"os"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/27.
// 执行CI步骤

type CIWorker struct {
	// Name 工程名称
	Name string
	// WorkDir 工作目录 默认为当前工程根目录
	WorkDir string
	// Hicd hicd配置数据
	Hicd model.HICD
}

func (c *CIWorker) Do() {
	dependenceResult, err := c.Dependence()
	if err != nil {
		log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s %s", "Dependence", dependenceResult, err.Error())}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s", "Dependence", dependenceResult)}, logrus.InfoLevel).Report()

	beforeResult, err := c.Before()
	if err != nil {
		log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s %s", "Before", beforeResult, err.Error())}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s", "Before", beforeResult)}, logrus.InfoLevel).Report()
}

// Dependence 执行构建前的依赖管理
// 如果need=true,则执行
func (c *CIWorker) Dependence() (string, error) {

	log.Output(model.DependenceModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("Work Dir %s", c.WorkDir)}, logrus.InfoLevel).Report()

	os.Chdir(c.WorkDir)

	result := ""
	if !c.Hicd.Dependence.Need {
		return result, nil
	}

	if len(c.Hicd.Dependence.Cmd) == 0 {
		return result, errors.New("Dependence Cmd Can Not Be Empty!")
	}

	return utils.CmdRun(c.Hicd.Dependence.Cmd)
}

func (c *CIWorker) Before() (string, error) {
	log.Output(model.BeforeModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("Work Dir %s", c.WorkDir)}, logrus.InfoLevel).Report()

	os.Chdir(c.WorkDir)

	if c.Hicd.Before.Skip {
		return " Skip [Before] Stage", nil
	}

	if len(c.Hicd.Before.Script) == 0 {
		return "", errors.New("Before Script Can Not Be Empty!")
	}

	return utils.CmdRun(c.Hicd.Before.Script)
}
