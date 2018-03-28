/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package worker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
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
		log.Output(model.DependenceModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s %s", "Dependence", err.Error())}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(model.DependenceModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s %s", "Dependence", dependenceResult)}, logrus.InfoLevel).Report()
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

	var out, stderr bytes.Buffer
	var cmd *exec.Cmd

	cmd = exec.Command("dep", "init", "-v")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		err = errors.New(fmt.Sprintf("%s\n%s", stderr.String(), err.Error()))
		return result, err
	}

	return out.String(), nil
	//logrus.Println(err)
	//logrus.Println(out.String())
	//logrus.Println(stderr.String())

	//cmd := exec.Command(c.Hicd.Dependence.Cmd[0], c.Hicd.Dependence.Cmd[1:]...)
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	return result, err
	//}
	//stderr, err := cmd.StderrPipe()
	//if err != nil {
	//	return result, err
	//}
	//
	//defer func() {
	//	stdout.Close()
	//	stderr.Close()
	//}()
	//
	//if err := cmd.Start(); err != nil {
	//	return result, err
	//}
	//
	//if err := cmd.Wait(); err != nil {
	//	return result, err
	//}
	//
	//var data []byte
	//_, err = stdout.Read(data)
	//if err != nil {
	//	return result, err
	//}
	//
	//result = string(data)
	//
	//_, err = stderr.Read(data)
	//if err != nil {
	//	return result, err
	//}
	//
	//if string(data) != "" {
	//	err = errors.New(string(data))
	//} else {
	//	err = nil
	//}

	return result, err
}
