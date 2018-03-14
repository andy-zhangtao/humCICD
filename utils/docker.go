/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/fsouza/go-dockerclient"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

func CreateContainer(opt model.BuildOpts) error {
	for _, o := range opt.DockerOpt {
		if err := buildContainer(opt.Client, o); err != nil {
			return err
		}
	}

	return nil
}

func buildContainer(cli *docker.Client, do model.DockerOpts) error {
	var env []string
	for key, value := range do.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	container, err := cli.CreateContainer(docker.CreateContainerOptions{
		Context: context.Background(),
		Config: &docker.Config{
			Image: do.Img,
		},
	})

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Name": do.Name, "ID": container.ID}).Info("BuildContainer")
	return cli.StartContainer(container.ID,&docker.HostConfig{})
}
