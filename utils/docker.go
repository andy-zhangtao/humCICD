/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"context"
	"fmt"
	"strings"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/model"
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

	container, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        do.Img,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env:          env,
		Cmd:          strings.Split(do.Cmd, " "),
	}, &container.HostConfig{
		AutoRemove: true,
	}, &network.NetworkingConfig{}, do.Name)

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Name": do.Name, "ID": container.ID}).Info("BuildContainer")

	return cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
}
