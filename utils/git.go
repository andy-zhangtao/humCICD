/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

// parseName 通过git地址解析工程名称
func ParseName(url string) (name string) {
	gitName := strings.Split(url, "/")
	name = strings.Split(gitName[len(gitName)-1], ".")[0]
	logrus.WithFields(log.Z().Fields(logrus.Fields{"GitAgent Will Clone": name})).Info(ModuleName)
	return
}

// ParsePath 通过git地址解析出clone后的路径
// 例如通过https://github.com/andy-zhangtao/humCICD.git提取 github.com/andy-zhangtao/humCICD
func ParsePath(url string) (path string) {
	if strings.HasPrefix(url, "https://") {
		path = url[len("https://") : len(url)-len(".git")]
	} else {
		path = url[len("http://") : len(url)-len(".git")]
	}

	return
}

// GetConfigure 调用API获取配置信息
// idOrName 工程ID或者名称 如果通过ID查询失败，就会通过名称查询
func GetConfigure(idOrName string) (*model.GitConfigure, error) {
	if os.Getenv(model.EnvDataAgent) == "" {
		return nil, log.Z().Error(fmt.Sprintf("[%s]Empty!", model.EnvDataAgent))
	}

	resp, err := getConfigureByID(idOrName)
	if err != nil {
		return nil, log.Z().Error(err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, log.Z().Error(err.Error())
		}

		var config model.GitConfigure

		err = json.Unmarshal(data, &config)
		if err != nil {
			logrus.WithFields(log.Z().Fields(logrus.Fields{"data": string(data)})).Info(ModuleName)
			return nil, log.Z().Error(err.Error())
		}

		return &config, nil
	case http.StatusNotFound:
		// 	通过ID查询失败
		resp, err = getConfigureByName(idOrName)
		if err != nil {
			return nil, log.Z().Error(err.Error())
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, log.Z().Error(err.Error())
		}

		var config model.GitConfigure

		err = json.Unmarshal(data, &config)
		if err != nil {
			logrus.WithFields(log.Z().Fields(logrus.Fields{"data": string(data)})).Info("GetConfigure")
			return nil, log.Z().Error(err.Error())
		}

		return &config, nil
	}

	return nil, log.Z().Error("API Invoke Error")
}

func DeleteConfigure(id string) (resp *http.Response, err error) {
	logrus.WithFields(log.Z().Fields(logrus.Fields{"API": os.Getenv(model.EnvDataAgent) + "/configure/" + id})).Info(ModuleName)
	client := &http.Client{}
	req, err1 := http.NewRequest("DELETE", os.Getenv(model.EnvDataAgent)+"/configure/"+id, nil)
	if err1 != nil {
		err = log.Z().Error(err1.Error())
		return
	}
	resp, err = client.Do(req)
	return
}

func getConfigureByID(id string) (resp *http.Response, err error) {
	api := os.Getenv(model.EnvDataAgent) + "/configure/" + id
	logrus.WithFields(log.Z().Fields(logrus.Fields{"API": api})).Info(ModuleName)
	resp, err = http.Get(api)
	return
}

func getConfigureByName(name string) (resp *http.Response, err error) {
	logrus.WithFields(log.Z().Fields(logrus.Fields{"API": os.Getenv(model.EnvDataAgent) + "/configure/name"})).Info(ModuleName)
	resp, err = http.Post(os.Getenv(model.EnvDataAgent)+"/configure/name", "", strings.NewReader(name))
	return
}
