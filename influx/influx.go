/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package influx

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/3/21.
const (
	LogDB = "hicd"
)

var influxCli client.Client

func init() {
	if os.Getenv(model.EnvInfluxDB) == "" {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Error": fmt.Sprintf("[%s]Empty", model.EnvInfluxDB)}, logrus.ErrorLevel)
		return
	}

	var err error
	influxCli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: os.Getenv(model.EnvInfluxDB),
	})

	if err != nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Connect InfluxDB Error": err}, logrus.ErrorLevel)
		return
	}

	_, version, err := influxCli.Ping(10 * time.Second)
	if err != nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Ping InfluxDB Error": err}, logrus.ErrorLevel)
		return
	}

	log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Connect InfluxDB Succ With Version": version}, logrus.InfoLevel)
}

// Insert 插入日志
// project 工程名称
// tags 标签MAP
// fields 属性MAP
func Insert(project string, tags, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  LogDB,
		Precision: "s",
	})
	if err != nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Create Batch Points Error": err}, logrus.ErrorLevel)
		return err
	}

	newTags := make(map[string]string)
	for k, v := range tags {
		if vv, ok := v.(string); ok {
			newTags[k] = vv
		}
	}

	pt, err := client.NewPoint(project, newTags, fields)
	if err != nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Create Points Error": err}, logrus.ErrorLevel)
		return err
	}

	bp.AddPoint(pt)

	return influxCli.Write(bp)
}

// Query 查询指定工程日志
// project 工程名称
func Query(project string) ([]model.RunLog, error) {
	var runLog []model.RunLog
	q := client.Query{
		Command:   fmt.Sprintf("select time, log from \"%s\"", project),
		Database:  LogDB,
		Precision: "s",
	}
	log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Query Str": q.Command}, logrus.InfoLevel)

	if response, err := influxCli.Query(q); err == nil && response.Error() == nil {
		res := response.Results

		for _, r := range res {
			for _, s := range r.Series {
				for _, v := range s.Values {
					rl := model.RunLog{}
					for i, vv := range v {
						if i == 0 {
							if vvv, ok := vv.(json.Number); ok {
								rl.Timestamp, _ = vvv.Int64()
							} else {
								logrus.Println(reflect.TypeOf(vv))
							}
						}
						if i == 1 {
							if vvv, ok := vv.(string); ok {
								rl.Message = vvv
							}
						}
					}
					runLog = append(runLog, rl)
				}
			}
		}
	} else {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Query Result Error": response.Err}, logrus.ErrorLevel)
		if response.Err != "" {
			return runLog, errors.New(fmt.Sprintf("%s %s", response.Err, err.Error()))
		}
		return runLog, err
	}

	return runLog, nil
}

// Destrory 销毁指定工程日志
func Destory(project string) error {
	q := client.Query{
		Command:   fmt.Sprintf("DROP MEASUREMENT \"%s\"", project),
		Database:  LogDB,
		Precision: "s",
	}

	log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Drop Str": q.Command}, logrus.InfoLevel)

	if response, err := influxCli.Query(q); err == nil && response.Error() == nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Destory Succ": true}, logrus.InfoLevel)
	} else {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Destroy Error": response.Err}, logrus.ErrorLevel)
		if response.Err != "" {
			return errors.New(fmt.Sprintf("%s %s", response.Err, err.Error()))
		}
		return err
	}

	return nil
}
