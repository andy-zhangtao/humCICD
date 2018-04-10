/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/andy-zhangtao/humCICD/db"
	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/nsqio/go-nsq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// dataAgent 用来将配置数据持久化到数据库当中. 并且提供查询API

var workerHome map[string]chan *nsq.Message
var workerChan chan *nsq.Message
/*buildAgent 从NSQ读取工程解析后的数据，然后执行构建任务*/
var producer *nsq.Producer

func nsqInit() {
	var errNum int
	var err error
	nsq_endpoint := os.Getenv(model.EnvNsqdEndpoint)
	if nsq_endpoint == "" {
		log.Output(model.DataAgent, "", logrus.Fields{"Env Empty": model.EnvNsqdEndpoint}, logrus.ErrorLevel).Report()
		// logrus.Error(fmt.Sprintf("[%s] Empty", model.EnvNsqdEndpoint))
		os.Exit(-1)
	}
	log.Output(model.DataAgent, "", logrus.Fields{"Connect NSQ": nsq_endpoint}, logrus.DebugLevel)
	for {
		producer, _ = nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
		err = producer.Ping()
		if err != nil {
			log.Output(model.DataAgent, "", logrus.Fields{"Ping Nsq Error": err}, logrus.ErrorLevel).Report()
			errNum ++
		}

		if err == nil {
			break
		}

		if errNum >= 20 {
			os.Exit(-1)
		}
		time.Sleep(time.Second * 5)
	}

	log.Output(model.DataAgent, "", logrus.Fields{"Connect Nsq Succes": producer.String()}, logrus.InfoLevel)
}

type DataAgent struct {
	Name        string
	NsqEndpoint string
}

func (this *DataAgent) HandleMessage(m *nsq.Message) error {
	logrus.WithFields(logrus.Fields{"HandleMessage": string(m.Body)}).Info(this.Name)
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *DataAgent) Run() {
	nsqInit()
	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(model.GitAgentTopic, this.Name, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": this.Name}).Error(this.Name)
		return
	}

	go func() {
		logrus.WithFields(logrus.Fields{"WorkChan": "Listen..."}).Info(this.Name)
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"BuildMsg": string(m.Body)}).Info(this.Name)
			msg := model.GitConfigure{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			go this.handleBuild(msg)

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&DataAgent{Name: this.Name}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{this.Name: "Listen...", "NSQ": this.NsqEndpoint}).Info(this.Name)
	<-r.StopChan
}

func (this *DataAgent) handleBuild(msg model.GitConfigure) {
	logrus.WithFields(logrus.Fields{"Name": msg.Name, "GitUrl": msg.GitUrl, "Configrue": msg.Configrue}).Info(this.Name)
	id, err := db.SaveConfig(&msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Name": msg.Name, "Save Mongo Error": err}).Error(this.Name)
	}

	err = producer.Publish(model.GitConfIDTopic, []byte(id))
	if err != nil {
		logrus.WithFields(logrus.Fields{"Name": msg.Name, "Send NSQ Error": err}).Error(this.Name)
	}
}

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'
		*/
		//"createTodo": &graphql.Field{
		//	Type:        model.ProjectType, // the return type for this field
		//	Description: "Create new todo",
		//	Args: graphql.FieldConfigArgument{
		//		"text": &graphql.ArgumentConfig{
		//			Type: graphql.NewNonNull(graphql.String),
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//
		//		// marshall and cast the argument value
		//		text, _ := params.Args["text"].(string)
		//
		//		// figure out new id
		//		newID := RandStringRunes(8)
		//
		//		// perform mutation operation here
		//		// for e.g. create a Todo and save to DB.
		//		newTodo := Todo{
		//			ID:   newID,
		//			Text: text,
		//			Done: false,
		//		}
		//
		//		TodoList = append(TodoList, newTodo)
		//
		//		// return the new Todo object that we supposedly save to DB
		//		// Note here that
		//		// - we are returning a `Todo` struct instance here
		//		// - we previously specified the return Type to be `todoType`
		//		// - `Todo` struct maps to `todoType`, as defined in `todoType` ObjectConfig`
		//		return newTodo, nil
		//	},
		//},
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:"a",done:true){id,text,done}}'
		*/
		"updateProject": &graphql.Field{
			Type:        model.ProjectType, // the return type for this field
			Description: "Update existing project, mark it activity or unactivity",
			Args: graphql.FieldConfigArgument{
				"status": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument value
				status, _ := params.Args["status"].(string)
				id, _ := params.Args["id"].(string)

				project := model.Project{
					Status: status,
				}

				logrus.WithFields(logrus.Fields{"id": id, "status": status}).Info(model.DataAgent)
				newProject, err := db.UpdateProject(id, project)
				if err != nil {
					return nil, err
				}

				return newProject, nil
			},
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/project?query={lastTodo{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		/*
		   curl -g 'http://localhost:8080/graphql?query={todo(id:"b"){id,text,done}}'
		*/
		"project": &graphql.Field{
			Type:        model.ProjectType,
			Description: "Get single project",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					// Search for el with id
					project, err := db.FindProjectByID(idQuery)
					if err != nil {
						return nil, err
					}
					p, err := model.Conver2Project(project)
					if err != nil {
						return nil, err
					}
					logrus.WithFields(logrus.Fields{"Project": p}).Info(model.DataAgent)
					return p, nil
				}
				return nil, nil
			},
		},

		"projects": &graphql.Field{
			Type:        graphql.NewList(model.ProjectType),
			Description: "Get All Projects",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				projects, err := db.FindAllProject()
				if err != nil {
					return nil, err
				}

				logrus.WithFields(logrus.Fields{"Projects": projects}).Debug(model.DataAgent)
				return projects, nil
			},
		},
		//"lastTodo": &graphql.Field{
		//	Type:        todoType,
		//	Description: "Last todo added",
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		return TodoList[len(TodoList)-1], nil
		//	},
		//},
		//
		///*
		//   curl -g 'http://localhost:8080/graphql?query={todoList{id,text,done}}'
		//*/
		//"todoList": &graphql.Field{
		//	Type:        graphql.NewList(todoType),
		//	Description: "List of todos",
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		return TodoList, nil
		//	},
		//},
	},
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query map[string]interface{}, schema graphql.Schema) *graphql.Result {

	params := graphql.Params{
		Schema:        schema,
		RequestString: query["query"].(string),
	}

	logrus.WithFields(logrus.Fields{"query": query}).Info(model.DataAgent)
	if query["variables"] != nil {
		params.VariableValues = query["variables"].(map[string]interface{})
	}

	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

// handleGraphQL 工程操作API
func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var g map[string]interface{}
	if r.Method == http.MethodGet {
		g = make(map[string]interface{})
		g["query"] = r.URL.Query().Get("query")
		result := executeQuery(g, schema)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(model.DataAgent)
		json.NewEncoder(w).Encode(result)
	}

	if r.Method == http.MethodPost {
		data, _ := ioutil.ReadAll(r.Body)
		logrus.WithFields(logrus.Fields{"body": string(data)}).Info(model.DataAgent)

		err := json.Unmarshal(data, &g)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
		}
		logrus.WithFields(logrus.Fields{"graph": g}).Info(model.DataAgent)
		result := executeQuery(g, schema)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(model.DataAgent)
		json.NewEncoder(w).Encode(result)
	}

}

func main() {
	bagent := DataAgent{
		Name:        model.DataAgent,
		NsqEndpoint: os.Getenv(model.EnvNsqdEndpoint),
	}

	go func() {
		router := mux.NewRouter()
		// 通过id查询git configure数据
		router.Path("/configure/{id:[0-9A-Za-z]+}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			id := mux.Vars(request)["id"]
			log.Output(model.DataAgent, "", logrus.Fields{"Find Configure ID": id}, logrus.InfoLevel)
			config, err := db.FindConfigByID(id)
			if err != nil {
				log.Output(model.DataAgent, "", logrus.Fields{"Find Configure Error": err}, logrus.ErrorLevel)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(&config)
		}).Name("GetConfigrue").Methods(http.MethodGet)

		// 通过name查询git configure数据
		router.Path("/configure/name").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			data, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Output(model.DataAgent, "", logrus.Fields{"Read Data Error": err}, logrus.ErrorLevel)
				return
			}

			log.Output(model.DataAgent, "", logrus.Fields{"Find Configure Name": string(data)}, logrus.InfoLevel)
			config, err := db.FindConfigByName(string(data))
			if err != nil {
				log.Output(model.DataAgent, "", logrus.Fields{"Find Configure Error": err}, logrus.ErrorLevel)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(&config)
		}).Name("GetConfigrue").Methods(http.MethodPost)

		// 通过id删除git configure数据
		router.Path("/configure/{id:[0-9A-Za-z]+}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			id := mux.Vars(request)["id"]
			log.Output(model.DataAgent, "", logrus.Fields{"Find Configure ID": id}, logrus.InfoLevel)
			err := db.DeleteConfigByID(id)
			if err != nil {
				log.Output(model.DataAgent, "", logrus.Fields{"Find Configure Error": err}, logrus.ErrorLevel)
				return
			}

		}).Name("GetConfigrue").Methods(http.MethodDelete)

		// GraphQLAPI接口,操作工程数据
		router.Path("/hicd").HandlerFunc(handleGraphQL)
		if err := http.ListenAndServe(":8000", cors.Default().Handler(router)); err != nil {
			log.Output(model.DataAgent, "", logrus.Fields{"Bind Port Error": err}, logrus.PanicLevel)
		}

	}()
	bagent.Run()

}
