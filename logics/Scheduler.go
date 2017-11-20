package logics

import (
	"github.com/hashicorp/nomad/api"
	"github.com/golang-plus/errors"
	"io/ioutil"
	"strconv"
	"encoding/base64"
	"os"
)

type Scheduler struct {
	nomad *api.Client
	jobsDir string
}

type ScheduleResponse struct {
	EvalId string `json:"eval_id"`
}

type UnscheduleResponse struct {
	EvalId string `json:"eval_id"`
}

func NewScheduler(nomadClient *api.Client, jobsDir string) (*Scheduler) {
	return &Scheduler{nomad: nomadClient, jobsDir: jobsDir}
}

func (s *Scheduler) Schedule(logic *Logic, logicVersion *LogicVersion) (*ScheduleResponse, error) {
	if logicVersion.Language == "javascript" {
		return s.scheduleJavascript(logic, logicVersion)
	} else {
		return nil, errors.New("Unsupported logic language " + logicVersion.Language)
	}
}

func (s *Scheduler) Unschedule(logicId string, logicVersionId int) (*UnscheduleResponse, error) {
	res, _, err := s.nomad.Jobs().Deregister(logicId  + "_" + strconv.Itoa(logicVersionId), true, &api.WriteOptions{})

	if err != nil {
		return nil, err
	} else {
		return &UnscheduleResponse{EvalId: res}, nil
	}
}

func (s *Scheduler) scheduleJavascript(logic *Logic, logicVersion *LogicVersion) (*ScheduleResponse, error) {
	data, err := base64.StdEncoding.DecodeString(logicVersion.Code)
	if err != nil {
		return nil, err
	}

	// -- create the application directory
	err = os.MkdirAll(s.jobsDir + "/" + logic.Id + "_" + strconv.Itoa(logicVersion.Version), 0777)
	if err != nil {
		return nil, err
	}

	// -- generate the application file
	err = ioutil.WriteFile(s.jobsDir + "/" + logic.Id + "_" + strconv.Itoa(logicVersion.Version) + "/app.js", data, 0777)
	if err != nil {
		return nil, err
	}

	task := api.NewTask(logic.Id + "_logic", "docker")

	task.Config = make(map[string]interface{})
	task.Config["image"] = "node:8"
	task.Config["command"] = "node"
	task.Config["args"] = []string{"app.js"}
	task.Config["volumes"] = []string{ s.jobsDir + "/" + logic.Id + "_" + strconv.Itoa(logicVersion.Version) + ":/usr/src/app" }
	task.Config["work_dir"] = "/usr/src/app"

	taskGroup := api.NewTaskGroup("logics", 1)
	taskGroup.Tasks = []*api.Task{task}

	job := api.NewServiceJob(logic.Id + "_" + strconv.Itoa(logicVersion.Version), logic.Id + " version " + strconv.Itoa(logicVersion.Version), "global", 1)
	job.Datacenters = []string{ "dc1" }
	job.TaskGroups = []*api.TaskGroup{taskGroup}

	resp, _, err := s.nomad.Jobs().Register(job, &api.WriteOptions{})

	if err != nil {
		return nil, err
	}

	return &ScheduleResponse{resp.EvalID}, nil
}