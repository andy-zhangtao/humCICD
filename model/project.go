/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

import (
	"github.com/graphql-go/graphql"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/6.

// Project 工程数据
type Project struct {
	ID     bson.ObjectId `json:"_id" bson:"_id"`
	Name   string        `json:"name"`
	Branch string        `json:"branch"`
	// Status 工程状态
	Status string `json:"status"`
}

// define custom Project ObjectType `projectType` for our Golang struct `Project`
// Note that
// - the fields in our projectType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var ProjectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "project",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if pro, ok := p.Source.(Project); ok {
					return pro.ID.Hex(), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"branch": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func Conver2Project(oldProject interface{}) (project Project, err error) {
	data, err := bson.Marshal(oldProject)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &project)
	if err != nil {
		return
	}

	return
}
