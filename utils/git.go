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

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

// parseName 通过git地址解析工程名称
func ParseName(url string) (name string) {
	gitName := strings.Split(url, "/")
	name = strings.Split(gitName[len(gitName)-1], ".")[0]
	fmt.Printf("GitAgent Will Clone [%s]\n", name)
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
func GetConfigure(idOrName string) (*model.GitConfigure, error) {
	if os.Getenv(model.EnvDataAgent) == "" {
		return nil, errors.New(fmt.Sprintf("[%s]Empty!", model.EnvDataAgent))
	}

	resp, err := getConfigureByID(idOrName)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var config model.GitConfigure

		err = json.Unmarshal(data, &config)
		if err != nil {
			logrus.WithFields(logrus.Fields{"data": string(data)}).Info("GetConfigure")
			return nil, err
		}

		return &config, nil
	case http.StatusNotFound:
		// 	通过ID查询失败
		resp, err = getConfigureByName(idOrName)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var config model.GitConfigure

		err = json.Unmarshal(data, &config)
		if err != nil {
			logrus.WithFields(logrus.Fields{"data": string(data)}).Info("GetConfigure")
			return nil, err
		}

		return &config, nil
	}

	return nil, errors.New("API Invoke Error")
}

func DeleteConfigure(id string) (resp *http.Response, err error) {
	logrus.WithFields(logrus.Fields{"API": os.Getenv(model.EnvDataAgent) + "/configure/" + id}).Info("DeleteConfigure")
	client := &http.Client{}
	req, err1 := http.NewRequest("DELETE", os.Getenv(model.EnvDataAgent) + "/configure/" + id, nil)
	if err1 != nil{
		err = err1
		return
	}
	resp, err = client.Do(req)
	return
}

func getConfigureByID(id string) (resp *http.Response, err error) {
	logrus.WithFields(logrus.Fields{"API": os.Getenv(model.EnvDataAgent) + "/configure/" + id}).Info("GetConfigure")
	resp, err = http.Get(os.Getenv(model.EnvDataAgent) + "/configure/" + id)
	return
}

func getConfigureByName(name string) (resp *http.Response, err error) {
	logrus.WithFields(logrus.Fields{"API": os.Getenv(model.EnvDataAgent) + "/configure/name"}).Info("GetConfigure")
	resp, err = http.Post(os.Getenv(model.EnvDataAgent)+"/configure/name", "", strings.NewReader(name))
	return
}
