/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

import "gopkg.in/mgo.v2/bson"

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/10.
type GitHubProject struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Refs        GitHubProject_refs `json:"refs"`
}

type GitHubProject_refs struct {
	TotalCount int                 `json:"totalcount"`
	Edges      []GitHubProjectNode `json:"edges"`
}

type GitHubProjectNode struct {
	Node GitHubProjectNode_node `json:"node"`
}

type GitHubProjectNode_node struct {
	Name string `json:"name"`
}

type GitHubSyncData struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	CreatedAt   string        `json:"created_at" bson:"createdAt"`
	UpdatedAt   string        `json:"updated_at" bson:"updatedAt"`
	Url         string        `json:"url" bson:"url"`
	Description string        `json:"description" bson:"description"`
	Branchs     []string      `json:"branchs" bson:"branchs"`
}
