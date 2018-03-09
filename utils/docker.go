/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"context"

	"docker.io/go-docker"
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
	//pb := make(map[docker.Port][]docker.PortBinding)
	//if len(do.Port) > 0 {
	//	for _, p := range do.Port {
	//		pb[docker.Port(fmt.Sprintf("%d/tcp", p))] = []docker.PortBinding{
	//			docker.PortBinding{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", p)},
	//		}
	//	}
	//}

	//logrus.WithFields(logrus.Fields{"Port": pb}).Info("BuildContainer")

	container, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        do.Img,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		AutoRemove: true,
	}, &network.NetworkingConfig{}, do.Name)
	//container, err := cli.CreateContainer(docker.CreateContainerOptions{
	//	Name: do.Name,
	//	Config: &docker.Config{
	//		Image:        do.Img,
	//		AttachStdout: true,
	//		AttachStdin:  true,
	//	},
	//	HostConfig: &docker.HostConfig{
	//		PortBindings: pb,
	//	},
	//	NetworkingConfig: &docker.NetworkingConfig{},
	//	Context:          context.Background(),
	//})

	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Name": do.Name, "ID": container.ID}).Info("BuildContainer")

	return nil
}
