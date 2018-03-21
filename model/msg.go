/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

// Write by zhangtao<ztao8607@gmail.com> . In 2018/2/24.

type TagEventMsg struct {
	Kind   string `json:"kind"`
	GitURL string `json:"git_url"`
	Tag    string `json:"tag"`
	Branch string `json:"branch"`
	Name   string `json:"name"`
}

type EventMsg struct {
	Kind  int         `json:"kind"`
	Msg   interface{} `json:"msg"`
	Email string      `json:"email"`
}

type PushEventMsg struct {
	GitURL string `json:"git_url"`
	Branch string `json:"branch"`
	Name   string `json:"name"`
}

type OutEventMsg struct {
	// Name 消息发生源
	Name string `json:"name"`
	// Proejct 工程名称
	Project string `json:"project"`
	// Out 消息内容
	Out string `json:"out"`
	// Result 当前消息状态 BuildSuc/BuildFailed
	Result int `json:"result"`
	// Time 消息时间戳
	Time string `json:"time"`
}
