package jobs

import (
	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"
	"github.org/dataprism/dataprism-kfunc/logics"
)

type JobsManager struct {
	consulClient *consul.KV
	nomadClient *nomad.Client
}

type JobListResult struct {
	Id string
	Name string
	Priority int
	Status string
	StatusDescription string
	SubmitTime int64
}

func NewManager(consulClient *consul.KV, nomadClient *nomad.Client) *JobsManager {
	return &JobsManager{
		consulClient: consulClient,
		nomadClient: nomadClient,
	}
}

func (m *JobsManager) List(prefix string) ([]JobListResult, error) {
	list, _, err := m.nomadClient.Jobs().PrefixList(prefix)

	if err != nil {
		return nil, err
	}

	res := make([]JobListResult, len(list))
	for _, v := range list {
		res = append(res, JobListResult{
			Id: v.ID,
			Name: v.Name,
			Priority: v.Priority,
			Status: v.Status,
			StatusDescription: v.StatusDescription,
			SubmitTime: v.SubmitTime,
		})
	}

	return res, nil
}

func (m *JobsManager) Get(id string) (*nomad.Job, error) {
	job, _, err := m.nomadClient.Jobs().Info(id, &nomad.QueryOptions{})

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (m *JobsManager) Schedule(logic *logics.Logic) (string, error) {
	job := nomad.NewServiceJob(logic.Id, logic.Id, "", 1)

	resp,_, err := m.nomadClient.Jobs().Register(job, &nomad.WriteOptions{})
	if err != nil {
		return "", err
	}

	return resp.EvalID, nil
}
