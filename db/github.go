/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/9.

const(
	ModuelName = "GitHubSync-DB-Service"
)

func SaveGitHubSync(s model.GitHubSyncData) (err error) {
	if s.ID == "" {
		s.ID = bson.NewObjectId()
	}
	return getGitHubSyncMongo().Insert(&s)
}

func DeleteGitHubSyncByID(id bson.ObjectId) (err error) {
	return getGitHubSyncMongo().RemoveId(id)
}

func GetAllGitHubSync() (sync []model.GitHubSyncData, err error) {
	err = getGitHubSyncMongo().Find(nil).All(&sync)
	return
}

func GetGitHubSyncByName(name string) (s model.GitHubSyncData, err error) {
	err = getGitHubSyncMongo().Find(bson.M{"name": name}).One(&s)
	return
}

func DeleteAllGitHubSync() (err error) {
	//_, err = getGitHubSyncMongo().RemoveAll(nil)
	//return
	_, err = bw.DeleteAll(new(model.GitHubSyncData))
	return err
}

func SaveALLGitHubSync(s []model.GitHubSyncData) (err error) {
	return bw.SaveAll(s)
}

//FindOneGitHubSync 查询符合条件的GitHubSync数据
//s 包含name属性的GithubSync数据
func FindOneGitHubSync(s *model.GitHubSyncData)(err error){
	logrus.WithFields(log.Z().Fields(logrus.Fields{"Receive Query Request":s})).Info(ModuelName)
	return bw.FindOne(s)
}