/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package db

import (
	"github.com/andy-zhangtao/humCICD/model"
	"gopkg.in/mgo.v2/bson"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/19.
// 保存参数相关数据

// SaveConfig 保存配置信息
func SaveConfig(config *model.GitConfigure) (string, error) {
	config.ID = bson.NewObjectId()
	err := getConfigureMongo().Insert(&config)
	return config.ID.Hex(), err
}

// FindConfigByID 根据_id返回配置信息
func FindConfigByID(id string) (configrue interface{}, err error) {
	err = getConfigureMongo().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&configrue)
	return
}

// DeleteConfigByID 根据_id删除配置信息
func DeleteConfigByID(id string) error {
	return getConfigureMongo().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
