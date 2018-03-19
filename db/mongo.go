/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/19.
var endpoint = os.Getenv(model.EnvMongo)
var username = os.Getenv(model.EnvMongoName)
var password = os.Getenv(model.EnvMongoPasswd)
var dbname = os.Getenv(model.EnvMongoDB)
var session *mgo.Session
var ModuleName = "Mongo-Init"

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

	logrus.WithFields(logrus.Fields{"Mongo Server": b.Version}).Info(ModuleName)
}

func getSession() *mgo.Session {
	return session.Clone()
}

func getConfigureMongo() *mgo.Collection {
	return getSession().DB(dbname).C(model.DefaultDBConf)
}
