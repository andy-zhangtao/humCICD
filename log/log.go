/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package log

import (
	"encoding/json"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/15.

type Log struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func Output(modelName string, fields logrus.Fields, level logrus.Level) *Log {
	switch level {
	case logrus.PanicLevel:
		logrus.WithFields(fields).Panic(modelName)
	case logrus.FatalLevel:
		logrus.WithFields(fields).Fatal(modelName)
	case logrus.ErrorLevel:
		logrus.WithFields(fields).Error(modelName)
	case logrus.WarnLevel:
		logrus.WithFields(fields).Warn(modelName)
	case logrus.InfoLevel:
		logrus.WithFields(fields).Info(modelName)
	case logrus.DebugLevel:
		logrus.WithFields(fields).Debug(modelName)
	}

	if _, ok := fields["msg"]; ok {
		content := fields["msg"]
		return &Log{Name: modelName, Content: content.(string)}
	}

	return &Log{Name: modelName, Content: "No Log"}
}

func (l *Log) Report() {
	output := model.OutEventMsg{
		Name:   l.Name,
		Out:    l.Content,
		Result: model.BuildSuc,
	}

	data, err := json.Marshal(output)
	if err != nil {
		Output(model.ReportTools, logrus.Fields{"Report Error": err}, logrus.ErrorLevel)
		return
	}

	utils.Reporter.Publish(model.HicdOutTopic, data)
}
