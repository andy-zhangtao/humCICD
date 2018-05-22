/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package hicdGraphql

import (
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/graphql-go/graphql"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/10.

var GitHubType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GitHubSync",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.ID.Hex(), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.Name, nil
				}
				return nil, nil
			},
		},
		"createAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updateAt": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.UpdatedAt, nil
				}
				return nil, nil
			},
		},
		"url": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.Url, nil
				}
				return nil, nil
			},
		},
		"description": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.Description, nil
				}
				return nil, nil
			},
		},
		"branch": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if g, ok := p.Source.(model.GitHubSyncData); ok {
					return g.Branchs, nil
				}
				return nil, nil
			},
		},
	},
})
