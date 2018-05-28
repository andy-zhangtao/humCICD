/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andy-zhangtao/humCICD/db"
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/sirupsen/logrus"
)

const (
	ModuleName = "GitHubSync-Service"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/9.

func SaveGitHubSync(sync []model.GitHubSyncData) (err error) {
	//for _, s := range sync {
	//	if err = db.SaveGitHubSync(s); err != nil {
	//		err = errors.New(fmt.Sprintf("Add GitHubSync Error [%s] sync [%v]", err.Error(), s))
	//		return
	//	}
	//}

	return db.SaveALLGitHubSync(sync)
}

func DeleGitHubSync(s model.GitHubSyncData) (err error) {
	if s.ID == "" {
		ts, err := db.GetGitHubSyncByName(s.Name)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}

			return err
		}
		s.ID = ts.ID
	}

	if err = db.DeleteGitHubSyncByID(s.ID); err != nil {
		err = errors.New(fmt.Sprintf("Delete GitHubSync Error [%s] Name[%s] ID[%s]", err.Error(), s.Name, s.ID))
	}

	return
}

func GetGitHubSync() (syncs []model.GitHubSyncData, err error) {
	return db.GetAllGitHubSync()
}

func RemoveAllGitHubSync() (err error) {
	if err = db.DeleteAllGitHubSync(); err != nil {
		err = errors.New(fmt.Sprintf("Remove ALl GitHubSync Error [%s]", err.Error()))
	}
	return
}

func GetSpecifyGitHubSync(sync *model.GitHubSyncData) (err error) {
	logrus.WithFields(log.Z().Fields(logrus.Fields{"Query github project Info": sync})).Info(ModuleName)
	err = db.FindOneGitHubSync(sync)
	logrus.WithFields(log.Z().Fields(logrus.Fields{"Return Github Project Sync Info": sync})).Info(ModuleName)
	return
}
