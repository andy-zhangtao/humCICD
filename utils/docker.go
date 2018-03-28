/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
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
		Config: &docker.Config{
			Env:          env,
			Cmd:          strings.Split(do.Cmd, " "),
			Image:        do.Img,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove: false,
			Binds:      do.Binds,
		},

		Context: context.Background(),
	})

	//container, err := cli.ContainerCreate(context.Background(), &container.Config{
	//	Image:        do.Img,
	//	Tty:          true,
	//	AttachStdout: true,
	//	AttachStderr: true,
	//	Env:          env,
	//	Cmd:          strings.Split(do.Cmd, " "),
	//}, &container.HostConfig{
	//	AutoRemove: true,
	//}, &network.NetworkingConfig{}, do.Name)

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Name": do.Name, "ID": container.ID}).Info("BuildContainer")
	//cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	return cli.StartContainer(container.ID, &docker.HostConfig{})
}
