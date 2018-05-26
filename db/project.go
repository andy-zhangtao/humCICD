/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"time"

	"github.com/andy-zhangtao/humCICD/model"
	"github.com/globalsign/mgo/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/6.

// FindAllProject 返回所有工程信息
func FindAllProject() (projects []model.Project, err error) {
	err = getProjectMongo().Find(nil).All(&projects)
	return
}

// FindProjectByID 根据_id返回工程信息
func FindProjectByID(id string) (project interface{}, err error) {
	err = getProjectMongo().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&project)
	return
}

// FindProjectByName 根据名称返回工程信息
func FindProjectByName(name string) (project interface{}, err error) {
	err = getProjectMongo().Find(bson.M{"name": name}).One(&project)
	return
}

// SaveConfig 保存配置信息
func SaveProject(project *model.Project) (string, error) {
	if project.ID == "" {
		project.ID = bson.NewObjectId()
	}

	if project.Time == "" {
		project.Time = time.Now().String()
	}
	err := getProjectMongo().Insert(&project)
	return project.ID.Hex(), err
}

// DeleteProjectByID 根据_id删除配置信息
func DeleteProjectByID(id string) error {
	return getProjectMongo().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func DeleteAllProject() error {
	_, err := getProjectMongo().RemoveAll(nil)
	return err
}
func UpdateProject(id string, project model.Project) (*model.Project, error) {
	tempProject, err := FindProjectByID(id)
	if err != nil {
		return nil, err
	}

	oldProject, err := model.Conver2Project(tempProject)
	if err != nil {
		return nil, err
	}

	if project.Name != "" {
		oldProject.Name = project.Name
	}

	if len(project.Branch) == 0 {
		oldProject.Branch = project.Branch
	}

	if project.Status != "" {
		oldProject.Status = project.Status
	}

	err = DeleteProjectByID(id)
	if err != nil {
		return nil, err
	}

	_, err = SaveProject(&oldProject)
	return &oldProject, err
}
