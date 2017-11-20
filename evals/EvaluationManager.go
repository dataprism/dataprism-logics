package evals

import (
	nomad "github.com/hashicorp/nomad/api"
	"strings"
)

type EvaluationManager struct {
	nomadClient *nomad.Client
}

type JobListResult struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Priority int `json:"priority"`
	Status string `json:"status"`
	StatusDescription string `json:"statusDescription"`
	SubmitTime int64 `json:"timestamp"`
}

func NewManager(nomadClient *nomad.Client) *EvaluationManager {
	return &EvaluationManager{
		nomadClient: nomadClient,
	}
}

func (m *EvaluationManager) Get(evaluationId string) (*Evaluation, error) {
	eval, _, err := m.nomadClient.Evaluations().Info(evaluationId, &nomad.QueryOptions{})

	if err != nil {
		return nil, err
	}

	if eval == nil {
		return nil, nil
	}

	return &Evaluation{
		Id:       eval.ID,
		Status:   eval.Status,
		Priority: eval.Priority,
		LogicId:  eval.JobID[:strings.LastIndex(eval.JobID, "_")],
		LogicVersion: eval.JobID[strings.LastIndex(eval.JobID, "_") + 1:],
	}, nil
}

func (m *EvaluationManager) Events(evaluationId string) (map[string]NodeAllocationState, error) {
	list, _, err := m.nomadClient.Evaluations().Allocations(evaluationId, &nomad.QueryOptions{})

	if err != nil {
		return nil, err
	}

	result := make(map[string]NodeAllocationState)
	for _, v := range list {
		res := NodeAllocationState{
			LogicId:  v.JobID[:strings.LastIndex(v.JobID, "_")],
			LogicVersion: v.JobID[strings.LastIndex(v.JobID, "_") + 1:],
			ActualStatus: v.ClientStatus,
			DesiredStatus: v.DesiredStatus,
			AllocationId: v.ID,
			NodeId: v.NodeID,
		}

		state := v.TaskStates[res.LogicId + "_logic"]

		for _, e := range state.Events {
			t := &Trace{Type: e.Type, Timestamp: e.Time}

			if e.RestartReason != "" { t.Message = e.RestartReason }
			if e.SetupError != "" { t.Message = e.SetupError }
			if e.DriverError != "" { t.Message = e.DriverError }
			if e.DriverMessage != "" { t.Message = e.DriverMessage }
			if e.KillError != "" { t.Message = e.KillError }
			if e.KillReason != "" { t.Message = e.KillReason }
			if e.DownloadError != "" { t.Message = e.DownloadError }
			if e.ValidationError != "" { t.Message = e.ValidationError }
			if e.VaultError != "" { t.Message = e.VaultError }
			if e.TaskSignalReason != "" { t.Message = e.TaskSignalReason }
			if e.Message != "" { t.Message = e.Message }

			res.Trace = append(res.Trace, t)
		}

		result[v.NodeID] = res
	}

	return result, nil
}
