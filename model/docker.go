/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

import "github.com/fsouza/go-dockerclient"

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

type BuildOpts struct {
	Client    *docker.Client
	DockerOpt []DockerOpts `json:"docker_opt"`
}

type DockerOpts struct {
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Port []int             `json:"port"`
	Env  map[string]string `json:"env"`
	Cmd  string            `json:"cmd"`
}
