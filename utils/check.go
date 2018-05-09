/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/andy-zhangtao/humCICD/model"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/8.

// CheckGitHubToken检查GitHub Token是否存在
func CheckGitHubToken() (err error) {
	if os.Getenv(model.Env_HICD_GitHub_Token) == "" {
		return errors.New(fmt.Sprintf("%s Empty!", model.Env_HICD_GitHub_Token))
	}

	return
}

func CheckMongo() (err error) {
	if os.Getenv(model.EnvMongo) == "" {
		return errors.New(fmt.Sprintf("%s Empty!", model.EnvMongo))
	}

	if os.Getenv(model.EnvMongoName) == "" {
		return errors.New(fmt.Sprintf("%s Empty!", model.EnvMongoName))
	}

	if os.Getenv(model.EnvMongoPasswd) == "" {
		return errors.New(fmt.Sprintf("%s Empty!", model.EnvMongoPasswd))
	}

	if os.Getenv(model.EnvMongoDB) == "" {
		return errors.New(fmt.Sprintf("%s Empty!", model.EnvMongoDB))
	}

	return
}
