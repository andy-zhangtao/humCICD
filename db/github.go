/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"github.com/andy-zhangtao/humCICD/model"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/9.

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
	_, err = getGitHubSyncMongo().RemoveAll(nil)
	return
}
