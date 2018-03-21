/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package influx

import (
	"fmt"
	"os"

	"github.com/andy-zhangtao/humCICD/log"
	"github.com/andy-zhangtao/humCICD/model"
	"github.com/influxdata/influxdb/client/v2"
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
		Addr: "http://localhost:8086",
	})

	if err != nil {
		log.Output(model.InfluxTools, model.DefualtEmptyProject, logrus.Fields{"Connect InfluxDB Error": err}, logrus.ErrorLevel)
		return
	}
}

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
