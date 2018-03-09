/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/humCICD/agents"
	"github.com/andy-zhangtao/humCICD/model"
)

const ModuleName = "Agent"
const BuildAgent = "buildagent"

func main() {
	if os.Getenv(model.EnvNsqdEndpoint) == ""{
		logrus.Panic(model.EnvNsqdEndpoint+" Empty")
	}

	logrus.WithFields(logrus.Fields{"HUM-AGENT":"START"}).Info(ModuleName)
	switch strings.ToLower(os.Getenv("HUM_AGENT")) {
	case BuildAgent:
		ba := agents.BuildAgent{Name:BuildAgent,NsqEndpoint:os.Getenv(model.EnvNsqdEndpoint)}
		ba.Run()
	}

}
