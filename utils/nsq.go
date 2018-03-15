/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"os"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/15.

var Reporter *nsq.Producer

func init() {
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(model.ReportTools)

	Reporter, err = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{"Connect Nsq Error": err,}).Panic(model.ReportTools)
	}

	err = Reporter.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Ping Nsq Error": err,}).Panic(model.ReportTools)
	}

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": Reporter.String(),}).Info(model.ReportTools)
}
