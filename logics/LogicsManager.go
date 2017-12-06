package logics

import (
	"context"
	"encoding/json"
	"github.com/dataprism/dataprism-commons/execute"
	"github.com/dataprism/dataprism-commons/core"
)

type LogicsManager struct {
	platform *core.Platform
}

func NewLogicsManager(platform *core.Platform) *LogicsManager {
	return &LogicsManager{platform: platform}
}

func (m *LogicsManager) ListLogics(ctx context.Context) ([]*Logic, error) {
	var result []*Logic

	pairs, err := m.platform.KV.List(ctx, "logics")
	if err != nil { return nil, err }

	for _, p := range pairs {
		var entity Logic
		if err = json.Unmarshal(p.Value, &entity); err != nil { return nil, err }
		result = append(result, &entity)
	}

	return result, err
}

func (m *LogicsManager) GetLogic(ctx context.Context, id string) (*Logic, error) {
	data, err := m.platform.KV.Get(ctx,"logics/" + id)
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

func (m *LogicsManager) GetLogicStatus(ctx context.Context, id string) (*execute.DataprismJobStatus, error) {
	return m.platform.Scheduler.GetJobStatus("logic", id)
}

func (m *LogicsManager) SetLogic(ctx context.Context, logic *Logic) (*Logic, error) {
	data, err := json.Marshal(logic)
	if err != nil {
		return nil, err
	}

	err = m.platform.KV.Set(ctx, "logics/" + logic.Id, data)
	if err != nil {
		return nil, err
	}

	return logic, nil
}

func (m *LogicsManager) RemoveLogic(ctx context.Context, id string) (error) {
	return m.platform.KV.Remove(ctx, "logics/" + id)
}