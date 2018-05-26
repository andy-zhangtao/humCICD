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
	// EnvTickTime 定时触发间隔时间
	EnvTickTime = "HICD_TICK"
	// EnvGitHubToken Github API Key
	EnvGitHubToken = "HICD_GITHUB_TOKEN"
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
	SpiderAgent  = "spiderAgent"
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
	BuildSuc = iota
	BuildFaild
)

const (
	GitImage = "vikings/gitagent:latest"
	GoImage  = "vikings/goagent:latest"
)

const (
	PushEventType = iota
	BranchTagEventType
)

const (
	DefaultDBName       = "hicd"
	DefaultDBConf       = "configure"
	DefaultProConf      = "project"
	DefaultGitHubSync   = "github_sync"
	DefualtEmptyProject = "SYSTEM_LOG"
	DefualtFinishFlag   = "Handler End"
	DependenceModule    = "Dependence"
	BeforeModule        = "Before"
	BuildModule         = "Build"
	TestModule          = "Test"
	ExceptionModule     = "Exception"
	AfterModule         = "After"
	WorkerModule        = "Hicd-Worker"
)

const (
	Env_HICD_GitHub_Token = "ENV_HICD_GITHUB_TOKEN"
	Env_HICD_GitHub_Name  = "ENV_HICD_GITHUB_NAME"
)

const (
	DB_GITHUB_SYNC = "github_sync"
)
