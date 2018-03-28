/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

import "gopkg.in/mgo.v2/bson"

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/7.
type Config struct {
	Title   string      `json:"title"`
	Version string      `json:"version"`
	Build   BuildConfig `json:"build"`
}

type BuildConfig struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type HicdConfigure struct {
	Version string     `json:"version"` // 当前配置文件版本
	Kind    string     `json:"kind"`    // 语言类型
	Depend  []HiDepend `json:"depend"`  // 依赖数据
	Dep     bool       `json:"dep"`     // 是否执行dep
}

type HiDepend struct {
	Kind string   `json:"kind"` // 依赖类型 docker, shell etc
	Exec string   `json:"exec"` // 执行命令内容, 如果是docker则是镜像名称, 如果是shell，则是shell 命令
	Args []string `json:"args"` // 命令参数
}

type GitConfigure struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name"`
	GitUrl    string        `json:"giturl"`
	Branch    string        `json:"branch"`
	Email     string        `json:"email" bson:"email"`
	Configrue HICD          `json:"configrue"`
}

type HICD struct {
	Language    string
	Dependence  Dependence
	Env         Env
	Before      Before
	Build       Build
	After       After
	Integration Integration
}

type Dependence struct {
	Need bool
	Cmd  []string
}

type Env struct {
	Skip bool
	Var  []map[string]string
}

type Before struct {
	Skip   bool
	Script []string
}

type Build struct {
	IsMake        bool
	Ispersistence bool
	Make          Make
	Cmd           Cmd
	Persistence   Persistence
	Test          Test
	Exception     Exception
}

type Exception struct {
	Cmd []string
}

type Make struct {
	Targets []string
}

type Cmd struct {
	Cmd []string
}

type Persistence struct {
	Path string
}

type Test struct {
	Cmd []string
}

type After struct {
	Usedocker  bool
	Dockerfile Dockerfile
	Var        []map[string]string
	Script     Script
}

type Dockerfile struct {
	Path string
}

type Script struct {
	Cmd []string
}
type Integration struct {
	Need bool
}

func Conver2GitConfigure(config interface{}) (configure GitConfigure, err error) {
	data, err := bson.Marshal(config)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &configure)
	if err != nil {
		return
	}

	return
}
