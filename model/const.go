/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/7.
const (
	EnvNsqdEndpoint = "HICD_NSQD_ENDPOINT"
	EnvEmailHost    = "HICD_EMAIL_HOST"
	EnvEmailUser    = "HICD_EMAIL_USER"
	EnvEmailPass    = "HICD_EMAIL_PASS"
	EnvEmailDest    = "HICD_EMAIL_DEST"
)

const (
	GitAgent     = "gitAgent"
	BuildAgent   = "buildAgent"
	EchoAgent    = "echoAgent"
	TrafficAgent = "trafficAgent"
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

const (
	GitImage = "vikings/gitagent:latest"
	GoImage  = "vikings/goagent:latest"
)
