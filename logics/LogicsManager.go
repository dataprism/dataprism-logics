package logics

import (
	"context"
	"encoding/json"
	"strconv"
	nomad "github.com/hashicorp/nomad/api"
	consul2 "github.com/dataprism/dataprism-commons/consul"
)

type LogicsManager struct {
	storage *consul2.ConsulStorage
	scheduler *Scheduler
	nomadClient *nomad.Client
}

func NewManager(storage *consul2.ConsulStorage, nomadClient *nomad.Client, scheduler *Scheduler) *LogicsManager {
	return &LogicsManager{storage: storage, nomadClient: nomadClient, scheduler:scheduler}
}

func (m *LogicsManager) ListLogics(ctx context.Context) ([]*Logic, error) {
	var result []*Logic

	pairs, err := m.storage.List(ctx, "logics")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity Logic
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }
		result = append(result, &entity)
	}

	return result, err
}

func (m *LogicsManager) ListLogicVersions(ctx context.Context, id string) ([]*LogicVersion, error) {
	var result []*LogicVersion

	pairs, err := m.storage.List(ctx, "logics/" + id + "/versions/")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity LogicVersion
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }

		result = append(result, &entity)
	}

	return result, err
}

func (m *LogicsManager) ListLogicVersionIds(ctx context.Context, id string) ([]int, error) {
	var result []int

	pairs, err := m.storage.List(ctx, "logics/" + id + "/versions/")
	if err != nil { return nil, err }

	for _, p := range pairs {
		v, err := strconv.Atoi(p.Key)
		if err != nil { return nil, err }

		result = append(result, v)
	}

	return result, err
}

func (m *LogicsManager) GetLogic(ctx context.Context, id string) (*Logic, error) {
	data, err := m.storage.Get(ctx,"logics/" + id + "/definition")
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

func (m *LogicsManager) GetLogicStatus(ctx context.Context, id string) (*LogicStatus, error) {
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
	list, err :=  m.ListLogicVersionIds(ctx, id)
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
	data, err := m.storage.Get(ctx, "logics/" + id + "/versions/" + strconv.Itoa(version))
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

	err = m.storage.Set(ctx, "logics/" + logic.Id + "/definition", data)
	if err != nil {
		return nil, err
	}

	return logic, nil
}

func (m *LogicsManager) SetLogicVersion(ctx context.Context, logicId string, logicVersion *LogicVersion) (*LogicVersion, error) {
	if logicVersion.Version == 0 {
		// -- determine the next logic version
		list, err :=  m.ListLogicVersionIds(ctx, logicId)
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

	err = m.storage.Set(ctx, "logics/" + logicId + "/versions/" + strconv.Itoa(logicVersion.Version), data)
	if err != nil {
		return nil, err
	}

	return logicVersion, nil
}

func (m *LogicsManager) RemoveLogic(ctx context.Context, id string) (error) {
	list, err := m.ListLogicVersionIds(ctx, id)

	for _, v := range list {
		if err = m.RemoveLogicVersion(ctx, id, v); err != nil {
			return err;
		}
	}

	return m.storage.Remove(ctx, "logics/" + id)
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
	return m.storage.Remove(ctx, "logics/" + id + "/versions/" + strconv.Itoa(version))
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