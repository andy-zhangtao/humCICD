/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/andy-zhangtao/humCICD/hicdGraphql"
	"github.com/andy-zhangtao/humCICD/service"
	"github.com/andy-zhangtao/humCICD/utils"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	ModelName = "HICD-GraphQL-Agent"
)

func main() {
	router := mux.NewRouter()
	router.Path("/api").HandlerFunc(handleDevExGraphQL)
	handler := cors.AllowAll().Handler(router)
	logrus.Fatal(http.ListenAndServe(":8000", handler))
}

func init() {
	if err := utils.CheckMongo(); err != nil {
		logrus.WithFields(logrus.Fields{"Check Mongo Error": err}).Error(ModelName)
		os.Exit(-1)
	}
}

var schemaDevex, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootDevexQuery,
	//Mutation: rootMutation,
})

var rootDevexQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"github": &graphql.Field{
			Type: graphql.NewList(hicdGraphql.GitHubType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if gt, err := service.GetGitHubSync(); err != nil {
					return nil, errors.New(fmt.Sprintf("Query GitHub Error [%s]", err.Error()))
				} else {
					return gt, nil
				}
			},
		},
	},
})

func handleDevExGraphQL(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var g map[string]interface{}
	if r.Method == http.MethodGet {
		g = make(map[string]interface{})
		g["query"] = r.URL.Query().Get("query")
		result := executeDevExQuery(g, schemaDevex)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(ModelName)
		json.NewEncoder(w).Encode(result)
	}

	if r.Method == http.MethodPost {
		data, _ := ioutil.ReadAll(r.Body)
		logrus.WithFields(logrus.Fields{"body": string(data)}).Info(ModelName)

		err := json.Unmarshal(data, &g)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
		}
		logrus.WithFields(logrus.Fields{"graph": g}).Info(ModelName)
		result := executeDevExQuery(g, schemaDevex)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(ModelName)
		json.NewEncoder(w).Encode(result)
	}
}

func executeDevExQuery(query map[string]interface{}, schema graphql.Schema) *graphql.Result {

	params := graphql.Params{
		Schema:        schema,
		RequestString: query["query"].(string),
	}

	if query["variables"] != nil {
		params.VariableValues = query["variables"].(map[string]interface{})
	}

	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		fmt.Println("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
