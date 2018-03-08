/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/7.
type HicdConfigure struct {
	Version string     `json:"version"` //当前配置文件版本
	Kind    string     `json:"kind"`    //语言类型
	Depend  []HiDepend `json:"depend"`  //依赖数据
}

type HiDepend struct {
	Kind string   `json:"kind"` //依赖类型 docker, shell etc
	Exec string   `json:"exec"` //执行命令内容, 如果是docker则是镜像名称, 如果是shell，则是shell 命令
	Args []string `json:"args"` //命令参数
}

type GitConfigure struct {
	Name      string        `json:"name"`
	GitUrl    string        `json:"git_url"`
	Branch    string        `json:"branch"`
	Configrue HicdConfigure `json:"configrue"`
}
