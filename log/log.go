/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package log

import (
	"encoding/json"
	"time"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/15.

type Log struct {
	Name    string `json:"name"`
	Proejct string `json:"proejct"`
	Result  int    `json:"result"`
	Content string `json:"content"`
}

func Output(modelName, project string, fields logrus.Fields, level logrus.Level) *Log {
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
		if level < logrus.WarnLevel {
			return &Log{Name: modelName, Proejct: project, Content: content.(string), Result: model.BuildFaild}
		}
		return &Log{Name: modelName, Proejct: project, Content: content.(string), Result: model.BuildSuc}

	}

	return &Log{Name: modelName, Proejct: project,}
}

func (l *Log) Report() {
	if l.Content == "" {
		// 如果消息为空, 不需要发送垃圾消息
		return
	}
	output := model.OutEventMsg{
		Name:    l.Name,
		Out:     l.Content,
		Project: l.Proejct,
		Result:  l.Result,
		Time:    time.Now().String(),
	}

	data, err := json.Marshal(output)
	if err != nil {
		Output(model.ReportTools, l.Proejct, logrus.Fields{"Report Error": err}, logrus.ErrorLevel)
		return
	}

	utils.Reporter.Publish(model.HicdOutTopic, data)
}
