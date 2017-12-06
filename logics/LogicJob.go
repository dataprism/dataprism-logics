package logics

import (
	"github.com/hashicorp/nomad/api"
	"github.com/dataprism/dataprism-commons/utils"
	"strings"
	"errors"
	"github.com/hashicorp/nomad/helper"
	"os"
	"io/ioutil"
	"encoding/base64"
	"github.com/dataprism/dataprism-commons/core"
)

type LogicJob struct {
	logic *Logic
	platform *core.Platform
}

func DefaultSyncJobResources() *api.Resources {
	return &api.Resources{
		CPU:      helper.IntToPtr(1500),
		MemoryMB: helper.IntToPtr(512),
		IOPS:     helper.IntToPtr(0),
	}
}

func NewSyncJob(logic *Logic, platform *core.Platform) LogicJob {
	return LogicJob{logic, platform }
}

func (s *LogicJob) ToJob() (*api.Job, error) {
	if s.logic.Language == "javascript" {
		return s.ToJavascriptJob(s.logic)
	} else {
		return nil, errors.New("Unsupported logic language " + s.logic.Language)
	}
}

func (s *LogicJob) ToJavascriptJob(logic *Logic) (*api.Job, error) {
	data, err := base64.StdEncoding.DecodeString(logic.Code)
	if err != nil { return nil, err }

	// -- create the application directory
	jobDir := s.platform.Settings.JobsDir + "/" + logic.Id
	if err = os.MkdirAll(jobDir, 0777); err != nil {
		return nil, err
	}

	// -- generate the application file
	if ioutil.WriteFile(jobDir + "/app.js", data, 0777); err != nil {
		return nil, err
	}

	nomadJobId := utils.ToNomadJobId("logic", s.logic.Id)

	task := api.NewTask(nomadJobId, "docker")

	task.Config = make(map[string]interface{})
	task.Config["image"] = "node:8"
	task.Config["command"] = "node"
	task.Config["args"] = []string{"app.js"}
	task.Config["volumes"] = []string{ jobDir + ":/usr/src/app" }
	task.Config["work_dir"] = "/usr/src/app"
	task.Env = make(map[string]string)
	task.Meta = make(map[string]string)
	task.Resources = DefaultSyncJobResources()

	taskGroup := api.NewTaskGroup("logics", 1)
	taskGroup.Tasks = []*api.Task{task}

	nomadJob := api.NewServiceJob(nomadJobId, strings.ToTitle(s.logic.Id), "global", 1)
	nomadJob.Datacenters = []string{ "aws" }
	nomadJob.TaskGroups = []*api.TaskGroup{taskGroup}

	return nomadJob, nil
}