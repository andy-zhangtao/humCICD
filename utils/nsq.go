/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"os"
	"time"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/15.

var Reporter *nsq.Producer

func init() {
	if os.Getenv(model.EnvNsqdEndpoint) == "" {
		logrus.WithFields(logrus.Fields{model.EnvNsqdEndpoint:"Empty!"}).Error(model.ReportTools)
		os.Exit(-1)
	}
	var err error
	var errNum int
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	logrus.WithFields(logrus.Fields{"Connect NSQ": nsq_endpoint,}).Info(model.ReportTools)

	for {
		Reporter, _ = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
		// if err != nil {
		// 	logrus.WithFields(logrus.Fields{"Connect Nsq Error": err,}).Error(model.ReportTools)
		// }

		err = Reporter.Ping()
		if err != nil {
			logrus.WithFields(logrus.Fields{"Ping Nsq Error": err,}).Error(model.ReportTools)
			errNum++
		}

		if err == nil {
			break
		}

		if errNum >= 20 {
			os.Exit(-1)
		}
		time.Sleep(time.Second * 5)
	}

	logrus.WithFields(logrus.Fields{"Connect Nsq Succes": Reporter.String(),}).Info(model.ReportTools)
}
