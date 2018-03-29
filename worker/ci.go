/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package worker

import (
	"fmt"
	"os"
	"strings"

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

	buildResult, err := c.Build()
	if err != nil {
		log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s %s", "Build", buildResult, err.Error())}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s", "Build", buildResult)}, logrus.InfoLevel).Report()

	afterResult, err := c.After()
	if err != nil {
		log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s %s", "After", afterResult, err.Error())}, logrus.ErrorLevel).Report()
		return
	}

	log.Output(model.WorkerModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s\n%s", "After", afterResult)}, logrus.InfoLevel).Report()
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

// Before 执行构建前的准备阶段
// 当skip=false执行此阶段脚本
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

// Build 触发构建
func (c *CIWorker) Build() (string, error) {
	log.Output(model.BuildModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("Work Dir %s", c.WorkDir)}, logrus.InfoLevel).Report()

	os.Chdir(c.WorkDir)

	testResult, err := testProject(c.Name, c.WorkDir, c.Hicd)
	if err != nil {
		out, _ := buildException(c.Name, c.WorkDir, c.Hicd)
		return fmt.Sprintf("%s\nException:\n%s", testResult, out), err
	}

	log.Output(model.TestModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s", testResult)}, logrus.InfoLevel).Report()

	buildResult, err := buildProject(c.Name, c.WorkDir, c.Hicd)
	if err != nil {
		out, _ := buildException(c.Name, c.WorkDir, c.Hicd)
		return fmt.Sprintf("%s\nException:\n%s", buildResult, out), err
	}

	log.Output(model.BuildModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%s", buildResult)}, logrus.InfoLevel).Report()
	return "", nil

}

// buildProject 构建工程
func buildProject(name, path string, conf model.HICD) (string, error) {

	hasMakefile := false

	_, err := os.Stat(path + "/Makefile")
	if err != nil {
		if os.IsNotExist(err) {
			hasMakefile = false
		}
	} else {
		hasMakefile = true
	}

	if conf.Build.IsMake {
		log.Output(model.BuildModule, name, logrus.Fields{"msg": "Use Makefile."}, logrus.InfoLevel).Report()
		//	使用makefile
		if hasMakefile {
			var makeCommand []string
			makeCommand = append(makeCommand, "make")
			log.Output(model.BuildModule, name, logrus.Fields{"msg": "Found Makefile"}, logrus.InfoLevel).Report()
			if len(conf.Build.Make.Targets) == 0 {
				return utils.CmdRun(makeCommand)
			}

			makeCommand = append(makeCommand, conf.Build.Make.Targets...)
			return utils.CmdRun(makeCommand)
		}

		log.Output(model.BuildModule, name, logrus.Fields{"msg": "Not Found Makefile. Use Buld Command"}, logrus.InfoLevel).Report()

		if len(conf.Build.Cmd.Cmd) == 0 {
			return utils.CmdRun([]string{"go", "build", "-v"})
		}
		return utils.CmdRun(conf.Build.Cmd.Cmd)
	}

	log.Output(model.BuildModule, name, logrus.Fields{"msg": "Don't Use Makefile. Use Buld Command"}, logrus.InfoLevel).Report()

	if len(conf.Build.Cmd.Cmd) == 0 {
		return utils.CmdRun([]string{"go", "build", "-v"})
	}
	return utils.CmdRun(conf.Build.Cmd.Cmd)
}

// testProject 构建前的单元测试
func testProject(name, path string, conf model.HICD) (string, error) {
	os.Chdir(path)

	if len(conf.Build.Test.Cmd) == 0 {
		log.Output(model.TestModule, name, logrus.Fields{"msg": "Use Default Unit Test Command"}, logrus.InfoLevel).Report()
		return utils.CmdRun([]string{"go", "test", "-v", "./..."})
	}

	var testCommand []string
	testCommand = append(testCommand, "go")
	testCommand = append(testCommand, conf.Build.Test.Cmd...)
	log.Output(model.TestModule, name, logrus.Fields{"msg": fmt.Sprintf("Use Unit Test Command [%v]", testCommand)}, logrus.InfoLevel).Report()
	return utils.CmdRun(testCommand)
}

// buildException 构建过程中的异常处理
func buildException(name, path string, conf model.HICD) (string, error) {
	log.Output(model.ExceptionModule, name, logrus.Fields{"msg": "Exception Appear"}, logrus.InfoLevel).Report()
	os.Chdir(path)

	if len(conf.Build.Exception.Cmd) != 0 {
		return utils.CmdRun(conf.Build.Exception.Cmd)
	}

	return "Excepiton Run Away", nil
}

// After 构建结束后的触发动作
func (c *CIWorker) After() (string, error) {
	os.Chdir(c.WorkDir)

	if c.Hicd.After.Usedocker {

		utils.CmdRun([]string{"apk", "add", "--update", "docker"})
		docker := "Dockerfile"
		var env []string
		if c.Hicd.After.Dockerfile.Path != "" {
			docker = c.Hicd.After.Dockerfile.Path
		}

		if len(c.Hicd.After.Var) > 0 {
			for _, v := range c.Hicd.After.Var {
				for key, value := range v {
					env = append(env, "--build-arg")
					env = append(env, fmt.Sprintf("%s=%s", key, value))
				}
			}
		}

		path := "."
		if strings.Contains(docker, "/") {
			//	dockerfile不在当前目录,需要解析目标目录
			ts := strings.Split(docker, "/")
			path = strings.Join(ts[:len(ts)-1], "/")

		}

		buildCommand := []string{"docker", "build", "-t", c.Hicd.After.Dockerfile.Name}
		buildCommand = append(buildCommand, env...)
		buildCommand = append(buildCommand, "-f", docker, path)

		log.Output(model.AfterModule, c.Name, logrus.Fields{"msg": fmt.Sprintf("%v", buildCommand)}, logrus.InfoLevel).Report()

		buildResult, err := utils.CmdRun(buildCommand)
		if err != nil {
			return buildResult, err
		}

		if len(c.Hicd.After.Script.Cmd) != 0 {
			return utils.CmdRun(c.Hicd.After.Script.Cmd)
		}

		return "Build After Stage End", nil
	}

	return "Skip After Stage", nil
}
