/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

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
