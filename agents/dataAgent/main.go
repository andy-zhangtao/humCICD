/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/andy-zhangtao/humCICD/db"
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/gorilla/mux"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// dataAgent 用来将配置数据持久化到数据库当中. 并且提供查询API
var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message
/*buildAgent 从NSQ读取工程解析后的数据，然后执行构建任务*/
var producer *nsq.Producer

func nsqInit() {
	var errNum int
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		log.Output(model.DataAgent, "", logrus.Fields{"Env Empty": model.EnvNsqdEndpoint}, logrus.ErrorLevel).Report()
		// logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	log.Output(model.DataAgent, "", logrus.Fields{"Connect NSQ": nsq_endpoint}, logrus.DebugLevel)
	for {
		producer, _ = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
		err = producer.Ping()
		if err != nil {
			log.Output(model.DataAgent, "", logrus.Fields{"Ping Nsq Error": err}, logrus.ErrorLevel).Report()
			errNum ++
		}

		if err == nil {
			break
		}

		if errNum >= 20 {
			os.Exit(-1)
		}
		time.Sleep(time.Second * 5)
	}

	log.Output(model.DataAgent, "", logrus.Fields{"Connect Nsq Succes": producer.String()}, logrus.InfoLevel)
}

type DataAgent struct {
	Name        string
	NsqEndpoint string
}

func (this *DataAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *DataAgent) Run() {
	nsqInit()
	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.GitAgentTopic, this.Name, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name}).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(logrus.Fields{"WorkChan": "Listen..."}).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
			msg := model.GitConfigure{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			go this.handleBuild(msg)

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&DataAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint}).Info(this.Name)
	<-r.StopChan
}

func (this *DataAgent) handleBuild(msg model.GitConfigure) {
	logrus.WithFields(logrus.Fields{"Name": msg.Name, "Configrue": msg.Configrue}).Info(this.Name)
	id, err := db.SaveConfig(&msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Name": msg.Name, "Save Mongo Error": err}).Error(this.Name)
	}

	err = producer.Publish(model.GitConfIDTopic, []byte(id))
	if err != nil {
		logrus.WithFields(logrus.Fields{"Name": msg.Name, "Send NSQ Error": err}).Error(this.Name)
	}
}

func main() {
	bagent := DataAgent{
		Name:        model.DataAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	go func() {
		router := mux.NewRouter()
		router.Path("/configure/{id:[0-9A-Za-z]+}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			id := mux.Vars(request)["id"]
			log.Output(model.DataAgent, "", logrus.Fields{"Find Configure ID":id}, logrus.InfoLevel)
			config, err := db.FindConfigByID(id)
			if err != nil{
				log.Output(model.DataAgent, "", logrus.Fields{"Find Configure Error": err}, logrus.ErrorLevel)
				return
			}
			//
			// data, err := json.Marshal(&config)
			// if err != nil{
			// 	log.Output(model.DataAgent, "", logrus.Fields{"Marshal Error": err}, logrus.ErrorLevel)
			// 	return
			// }

			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(&config)
		}).Name("GetConfigrue").Methods(http.MethodGet)

		if err := http.ListenAndServe(":8000", router); err != nil {
			log.Output(model.DataAgent, "", logrus.Fields{"Bind Port Error": err}, logrus.PanicLevel)
		}
	}()
	bagent.Run()

}
