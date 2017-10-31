package logics

import (
	"context"
	"github.com/hashicorp/consul/api"
	"encoding/json"
)

type LogicsManager struct {
	client *api.KV
}

func NewManager(client *api.KV) *LogicsManager {
	return &LogicsManager{client: client}
}

func (m *LogicsManager) ListLogics(ctx context.Context) ([]string, error) {
	pairs, _, err := m.client.Keys("logics/", "/", &api.QueryOptions{})
	if err != nil {
		return nil, err
	} else {
		return pairs, nil
	}
}

func (m *LogicsManager) GetLogic(ctx context.Context, id string) (*Logic, error) {
	data, _, err := m.client.Get("logics/" + id, &api.QueryOptions{})
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

func (m *LogicsManager) SetLogic(ctx context.Context, logic *Logic) (*Logic, error) {
	data, err := json.Marshal(logic)
	if err != nil {
		return nil, err
	}

	pair := &api.KVPair{Key: "logics/" + logic.Id, Value: data}

	_, err = m.client.Put(pair, &api.WriteOptions{})
	if err != nil {
		return nil, err
	}

	return logic, nil
}

func (m *LogicsManager) RemoveLogic(ctx context.Context, id string) (error) {
	_, err := m.client.Delete("logics/" + id, &api.WriteOptions{})

	return err
}