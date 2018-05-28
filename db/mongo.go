/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/andy-zhangtao/bwidow"
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/19.
var endpoint = os.Getenv(model.EnvMongo)
var username = os.Getenv(model.EnvMongoName)
var password = os.Getenv(model.EnvMongoPasswd)
var dbname = os.Getenv(model.EnvMongoDB)
var session *mgo.Session
var ModuleName = "Mongo-Init"
var bw *bwidow.BW

func check() error {
	if endpoint == "" {
		return errors.New(fmt.Sprintf("[%s] Not Found", model.EnvMongo))
	}

	return nil
}

func init() {
	logrus.Println("=====Connect Mongo=====")
	err := check()
	if err != nil {
		logrus.Panic(err)
	}

	if dbname == "" {
		dbname = model.DefaultDBName
	}

	if username != "" || password != "" {
		dialInfo := &mgo.DialInfo{
			Addrs:    []string{endpoint},
			Database: dbname,
			Username: username,
			Password: password,
		}

		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
	} else {
		session, err = mgo.Dial(endpoint)
	}
	b, err := session.BuildInfo()
	if err != nil {
		panic(err)
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{"Mongo Server": b.Version})).Info(ModuleName)

	//	Use Black Widow
	bw = bwidow.GetWidow()
	err = bw.Driver(bwidow.DRIVER_MONGO).Error()
	if err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"BWidow Driver Init Error": err})).Errorln(ModuleName)
		return
	}

	bw.Map(model.GitHubSyncData{}, model.DB_GITHUB_SYNC)
	bw.Map(model.GitConfigure{}, model.DB_PROJECT_CONFIGURE)

	if err = bw.CheckIndex(new(model.GitHubSyncData)).CheckIndex(new(model.GitConfigure)).Error(); err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"BWidow Check Index Error": err.Error()})).Error(ModuleName)
	} else {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"BWidow Init Sucess Version": bw.Version()})).Info(ModuleName)
	}
}

func getSession() *mgo.Session {
	return session.Clone()
}

func getConfigureMongo() *mgo.Collection {
	return getSession().DB(dbname).C(model.DefaultDBConf)
}

func getProjectMongo() *mgo.Collection {
	return getSession().DB(dbname).C(model.DefaultProConf)
}

func getGitHubSyncMongo() *mgo.Collection {
	return getSession().DB(dbname).C(model.DefaultGitHubSync)
}
