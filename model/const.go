/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/7.
const (
	EnvNsqdEndpoint = "HICD_NSQD_ENDPOINT"
)

const (
	GitAgent   = "gitAgent"
	BuildAgent = "buildAgent"
	EchoAgent  = "echoAgent"
)

const (
	GitAgentTopic = "HICD_GitAgent"
	/*HicdErrTopic 错误信息*/
	HicdErrTopic = "HICD_Error"
	/*HicdOutTopic 正常输出信息*/
	HicdOutTopic = "HICD_Output"
	TAGQUEUE     = "HUM_GIT_TAG"
)

const (
	BuildSuc   = iota
	BuildFaild
)
