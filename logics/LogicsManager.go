package logics

import (
	"context"
	consul "github.com/hashicorp/consul/api"
	"encoding/json"
	"strings"
	"strconv"
	nomad "github.com/hashicorp/nomad/api"
)

type LogicsManager struct {
	client *consul.KV
	scheduler *Scheduler
	nomadClient *nomad.Client
}

func NewManager(client *consul.KV, nomadClient *nomad.Client, scheduler *Scheduler) *LogicsManager {
	return &LogicsManager{client: client, nomadClient: nomadClient, scheduler:scheduler}
}

func (m *LogicsManager) ListLogics(ctx context.Context) ([]string, error) {
	pairs, _, err := m.client.Keys("logics/", "/", &consul.QueryOptions{})
	if err != nil {
		return nil, err
	} else {
		if pairs == nil {
			return []string{}, nil
		}

		var res []string
		for _, p := range pairs {
			idx := strings.Index(p, "/")

			if idx == -1 {
				continue
			}

			idx2 := strings.Index(p[idx + 1:], "/")

			if idx2 == -1 {
				idx2 = len(p[idx + 1:])
			}

			res = append(res, p[idx + 1:][:idx2])
		}

		return res, nil
	}
}

func (m *LogicsManager) ListLogicVersions(ctx context.Context, id string) ([]int, error) {
	pairs, _, err := m.client.Keys("logics/" + id + "/versions/", "/", &consul.QueryOptions{})
	if err != nil {
		return nil, err
	} else {
		if pairs == nil {
			return []int{}, nil
		}

		var res []int
		for _, p := range pairs {
			idx := strings.LastIndex(p, "/")

			if idx == -1 {
				continue
			}

			i, err := strconv.Atoi(p[idx + 1:])

			if err == nil {
				res = append(res, i)
			}
		}

		return res, nil
	}
}

func (m *LogicsManager) GetLogic(ctx context.Context, id string) (*Logic, error) {
	data, _, err := m.client.Get("logics/" + id + "/definition", &consul.QueryOptions{})
	if err != nil {
		return nil, err
	}

	var logic Logic
	err = json.Unmarshal(data.Value, &logic)
	if err != nil {
		return nil, err
	}

	return &logic, nil
}

func (m *LogicsManager) GetLogicStatus(ctx context.Context, logicId string) (*LogicStatus, error) {
	list, err :=  m.ListLogicVersions(ctx, logicId)
	if err != nil {
		return nil, err
	}

	latest := 0;
	for _, e := range list {
		if latest < e {
			latest = e
		}
	}

	id := logicId + "_" + strconv.Itoa(latest)

	summary, _, err := m.nomadClient.Jobs().Summary(id, &nomad.QueryOptions{})
	if err != nil {
		return nil, err
	}

	tgSummary := summary.Summary["logics"]

	return &LogicStatus{
		Complete: tgSummary.Complete,
		Failed: tgSummary.Failed,
		Lost: tgSummary.Lost,
		Queued: tgSummary.Queued,
		Running: tgSummary.Running,
		Starting: tgSummary.Starting,
	}, nil
}

func (m *LogicsManager) GetLatestLogicVersion(ctx context.Context, id string) (*LogicVersion, error) {
	// -- determine the next logic version
	list, err :=  m.ListLogicVersions(ctx, id)
	if err != nil {
		return nil, err
	}

	latest := 0;
	for _, e := range list {
		if latest < e {
			latest = e
		}
	}

	return m.GetLogicVersion(ctx, id, latest)
}

func (m *LogicsManager) GetLogicVersion(ctx context.Context, id string, version int) (*LogicVersion, error) {
	data, _, err := m.client.Get("logics/" + id + "/versions/" + strconv.Itoa(version), &consul.QueryOptions{})
	if err != nil {
		return nil, err
	}

	if data != nil {
		var logicVersion LogicVersion
		err = json.Unmarshal(data.Value, &logicVersion)
		if err != nil {
			return nil, err
		}

		return &logicVersion, nil
	} else {
		return nil, nil
	}
}

func (m *LogicsManager) SetLogic(ctx context.Context, logic *Logic) (*Logic, error) {
	data, err := json.Marshal(logic)
	if err != nil {
		return nil, err
	}

	pair := &consul.KVPair{Key: "logics/" + logic.Id + "/definition", Value: data}

	_, err = m.client.Put(pair, &consul.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return logic, nil
}

func (m *LogicsManager) SetLogicVersion(ctx context.Context, logicId string, logicVersion *LogicVersion) (*LogicVersion, error) {
	if logicVersion.Version == 0 {
		// -- determine the next logic version
		list, err :=  m.ListLogicVersions(ctx, logicId)
		if err != nil {
			return nil, err
		}

		latest := 0;
		for _, e := range list {
			if latest < e {
				latest = e
			}
		}

		logicVersion.Version = latest + 1;
	}

	data, err := json.Marshal(logicVersion)
	if err != nil {
		return nil, err
	}

	pair := &consul.KVPair{Key: "logics/" + logicId + "/versions/" + strconv.Itoa(logicVersion.Version), Value: data}

	_, err = m.client.Put(pair, &consul.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return logicVersion, nil
}

func (m *LogicsManager) RemoveLogic(ctx context.Context, id string) (error) {
	list, err := m.ListLogicVersions(ctx, id)

	for _, v := range list {
		if err = m.RemoveLogicVersion(ctx, id, v); err != nil {
			return err;
		}
	}

	_, err = m.client.DeleteTree("logics/" + id, &consul.WriteOptions{})

	return err
}

func (m *LogicsManager) RemoveLogicVersion(ctx context.Context, id string, version int) (error) {
	list, _, err := m.nomadClient.Jobs().PrefixList(id + "_" + strconv.Itoa(version))
	if err != nil {
		return err
	}

	if list != nil && len(list) == 1 {
		if _, err = m.Unschedule(ctx, id, version); err != nil {
			return err
		}
	}
	_, err = m.client.Delete("logics/" + id + "/versions/" + strconv.Itoa(version), &consul.WriteOptions{})

	return err
}

func (m *LogicsManager) Schedule(ctx context.Context, id string, version int) (*ScheduleResponse, error) {
	logic, err := m.GetLogic(ctx, id)
	if err != nil {
		return nil, err
	}

	logicVersion, err := m.GetLogicVersion(ctx, id, version)
	if err != nil {
		return nil, err
	}

	return m.scheduler.Schedule(logic, logicVersion)
}

func (m *LogicsManager) Unschedule(ctx context.Context, id string, version int) (*UnscheduleResponse, error) {
	return m.scheduler.Unschedule(id, version)
}