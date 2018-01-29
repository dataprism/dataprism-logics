package logics

import (
	"context"
	"github.com/dataprism/dataprism-commons/execute"
	"github.com/dataprism/dataprism-commons/core"
)

type ExecutionManager struct {
	platform *core.Platform
	logicsManager *LogicsManager
}

func NewExecutionManager(platform *core.Platform, logicsManager *LogicsManager) *ExecutionManager {
	return &ExecutionManager{platform, logicsManager}
}

func (m *ExecutionManager) Deploy(ctx context.Context, id string) (*execute.ScheduleResponse, error) {
	// -- get the logic
	logic, err := m.logicsManager.GetLogic(ctx, id)
	if err != nil {
		return nil, err
	}

	// -- create the job for the link
	job := NewLogicJob(logic, m.platform)

	// -- schedule the job
	return m.platform.Scheduler.Schedule(job)
}

func (m *ExecutionManager) Undeploy(ctx context.Context, id string) (*execute.UnscheduleResponse, error) {
	// -- schedule the job
	return m.platform.Scheduler.Unschedule("logic", id)
}