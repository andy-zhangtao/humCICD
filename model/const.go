/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/7.
const (
	EnvNsqdEndpoint = "HICD_NSQD_ENDPOINT"
	EnvEmailHost    = "HICD_EMAIL_HOST"
	EnvEmailUser    = "HICD_EMAIL_USER"
	EnvEmailPass    = "HICD_EMAIL_PASS"
	EnvEmailDest    = "HICD_EMAIL_DEST"
	EnvMongo        = "HICD_MONGO_ENDPOINT"
	EnvMongoName    = "HICD_MONGO_NAME"
	EnvMongoPasswd  = "HICD_MONGO_PASSWD"
	EnvMongoDB      = "HICD_MONGO_DB"
	// EnvDataAgent dataAgent Endpoint
	EnvDataAgent = "HICD_DATA_AGENT"
	// EnvInfluxDB InfluxDB地址
	EnvInfluxDB = "HICD_INFLUX_DB"
)

const (
	GitAgent     = "gitAgent"
	BuildAgent   = "buildAgent"
	EchoAgent    = "echoAgent"
	DataAgent    = "dataAgent"
	GoAgent      = "goagent"
	TrafficAgent = "trafficAgent"
	ReportTools  = "report"
	InfluxTools  = "Influx"
)

const (
	GitAgentTopic = "HICD_GitAgent"
	/*HicdErrTopic 错误信息*/
	HicdErrTopic = "HICD_Error"
	/*HicdOutTopic 正常输出信息*/
	HicdOutTopic = "HICD_Output"
	TAGQUEUE     = "HUM_GIT_TAG"
	// GitConfIDTopic Git配置信息ID
	GitConfIDTopic = "HICD_Git_Config_Topic"
)

const (
	BuildSuc   = iota
	BuildFaild
)

const (
	GitImage = "vikings/gitagent:latest"
	GoImage  = "vikings/goagent:latest"
)

const (
	PushEventType      = iota
	BranchTagEventType
)

const (
	DefaultDBName       = "hicd"
	DefaultDBConf       = "configure"
	DefualtEmptyProject = "SYSTEM_LOG"
	DefualtFinishFlag   = "Handler End"
	DependenceModule    = "Dependence"
	BeforeModule        = "Before"
	WorkerModule        = "Hicd-Worker"
)
