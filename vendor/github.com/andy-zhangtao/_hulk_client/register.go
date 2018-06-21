package _hulk_client

import (
	"github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
	"os"
	"context"
	"fmt"
)

// Register 服务注册.
// Hulk支持服务注册, 每个通过Hulk获取配置文件的服务，可以将自己注册到Hulk当中
// 注册信息包括: name:服务名称 version:使用版本 resume:服务简介
// Hulk会记录以上三个信息, 同时记录服务IP
func Register(name, version, resume string) (err error) {
	variables := map[string]interface{}{
		"name":    graphql.String(name),
		"version": graphql.String(version),
		"resume":  graphql.String(resume),
	}

	var query struct {
		AddRegister struct {
			Name    graphql.String
		} `graphql:"addRegister(name: $name, version: $version, resume: $resume)"`
	}

	logrus.WithFields(logrus.Fields{"variables": variables}).Info(HULK_GO_SDK)

	client := graphql.NewClient(os.Getenv(ENDPOINT), nil)
	err = client.Mutate(context.Background(), &query, variables)
	if err != nil {
		logrus.Error(fmt.Sprintf("Service Register Error [%s]", err))
	}

	return
}

// UnRegister 服务卸载
// 当服务确定下线或者暂时离线时，调用此函数来通知Hulk服务下线
// name:服务名称 version:使用版本
func UnRegister(name, version string)(err error){
	variables := map[string]interface{}{
		"name":    graphql.String(name),
		"version": graphql.String(version),
	}

	var query struct {
		AddRegister struct {
			Name    graphql.String
		} `graphql:"deleteRegister(name: $name, version: $version)"`
	}

	logrus.WithFields(logrus.Fields{"variables": variables}).Info(HULK_GO_SDK)

	client := graphql.NewClient(os.Getenv(ENDPOINT), nil)
	err = client.Mutate(context.Background(), &query, variables)
	if err != nil {
		logrus.Error(fmt.Sprintf("Service Register Error [%s]", err))
	}

	return
}