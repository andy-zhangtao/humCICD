/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/andy-zhangtao/humCICD/service"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/globalsign/mgo/bson"
	"github.com/shurcooL/graphql"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	ModelName = "HICD-GitHub-Agent"
)

func main() {
	//定时查询GitHub 工程
	go syncGitHubInTime()

	<-make(chan int)

}

func init() {
	if err := utils.CheckGitHubToken(); err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Check GitHub Token Error": err})).Error(ModelName)
		os.Exit(-1)
	}

	if err := utils.CheckMongo(); err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Check Mongo Error": err})).Error(ModelName)
		os.Exit(-1)
	}


}

func syncGitHubInTime() {
	if err := syncGitHub(); err != nil {
		logrus.WithFields(log.Z().Fields(logrus.Fields{"Sync Repository Error": err})).Error(ModelName)
	}
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		fmt.Printf("下次采集时间为[%s]\n", next.Format("200601021504"))

		select {
		case <-t.C:
			err := syncGitHub()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func syncGitHub() (err error) {

	var syncData []model.GitHubSyncData

	auth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(model.Env_HICD_GitHub_Token)},
	)

	var query struct {
		Viewer struct {
			Repositories struct {
				Edges []struct {
					Cursor string
					Node struct {
						Name        string
						CreatedAt   string
						UpdatedAt   string
						Url         string
						Description string
					}
				}
			} `graphql:"repositories(first:50)""`
		}
	}

	httpClient := oauth2.NewClient(context.Background(), auth)
	client := graphql.NewClient("https://api.github.com/graphql", httpClient)
	err = client.Query(context.Background(), &query, nil)
	if err != nil {
		return
	}
	//logrus.WithFields(logrus.Fields{"total": len(query.Viewer.Repositories.Edges)}).Info(ModelName)
	//logrus.WithFields(logrus.Fields{"viewer": query.Viewer.Repositories.Edges[len(query.Viewer.Repositories.Edges)-1].Cursor}).Info(ModelName)

	for _, g := range query.Viewer.Repositories.Edges {
		branchs, err := syncBranchs(g.Node.Name)
		if err != nil {
			logrus.WithFields(log.Z().Fields(logrus.Fields{"Query Branch Error": err.Error(), "Repository": g.Node.Name})).Error(ModelName)
			continue
		}

		syncData = append(syncData, model.GitHubSyncData{
			ID:          bson.NewObjectId(),
			Name:        g.Node.Name,
			CreatedAt:   g.Node.CreatedAt,
			UpdatedAt:   g.Node.UpdatedAt,
			Url:         g.Node.Url,
			Description: g.Node.Description,
			Branchs:     branchs,
		})
	}
	cursor := query.Viewer.Repositories.Edges[len(query.Viewer.Repositories.Edges)-1].Cursor
	for {

		var queryWithCursor struct {
			Viewer struct {
				Repositories struct {
					Edges []struct {
						Cursor string
						Node struct {
							Name        string
							CreatedAt   string
							UpdatedAt   string
							Url         string
							Description string
						}
					}
				} `graphql:"repositories(first:50, after:$after)""`
			}
		}

		variables := map[string]interface{}{
			"after": graphql.NewString(graphql.String(cursor)),
		}

		err = client.Query(context.Background(), &queryWithCursor, variables)
		if err != nil {
			return
		}

		if len(queryWithCursor.Viewer.Repositories.Edges) == 0 {
			break
		}
		//logrus.WithFields(logrus.Fields{"size": len(queryWithCursor.Viewer.Repositories.Edges)}).Info(ModelName)

		cursor = queryWithCursor.Viewer.Repositories.Edges[len(queryWithCursor.Viewer.Repositories.Edges)-1].Cursor
		for _, g := range queryWithCursor.Viewer.Repositories.Edges {
			branchs, err := syncBranchs(g.Node.Name)
			if err != nil {
				logrus.WithFields(log.Z().Fields(logrus.Fields{"Query Branch Error": err.Error(), "Repository": g.Node.Name})).Error(ModelName)
				continue
			}

			syncData = append(syncData, model.GitHubSyncData{
				ID:          bson.NewObjectId(),
				Name:        g.Node.Name,
				CreatedAt:   g.Node.CreatedAt,
				UpdatedAt:   g.Node.UpdatedAt,
				Url:         g.Node.Url,
				Description: g.Node.Description,
				Branchs:     branchs,
			})
		}
	}

	logrus.WithFields(log.Z().Fields(logrus.Fields{"Sync Data": len(syncData)})).Info(ModelName)

	if err = service.RemoveAllGitHubSync(); err != nil {
		return
	}
	return service.SaveGitHubSync(syncData)
}

func syncBranchs(name string) (branchs []string, err error) {
	var query struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Name string `json:"name"`
					}
				}
			} `graphql:"refs(refPrefix:$refPrefix,first:100)"`
		} `graphql:"repository(owner:$owner,name:$name)"`
	}

	variables := map[string]interface{}{
		"owner":     graphql.String(os.Getenv(model.Env_HICD_GitHub_Name)),
		"name":      graphql.String(name),
		"refPrefix": graphql.String("refs/heads/"),
	}

	auth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(model.Env_HICD_GitHub_Token)},
	)

	httpClient := oauth2.NewClient(context.Background(), auth)
	client := graphql.NewClient("https://api.github.com/graphql", httpClient)
	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	for _, e := range query.Repository.Refs.Edges {
		branchs = append(branchs, e.Node.Name)
	}

	return
}
