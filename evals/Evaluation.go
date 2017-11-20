package evals

type Evaluation struct {
	Id       string `json:"id"`
	LogicId  string `json:"logic_id"`
	LogicVersion string `json:"logic_version"`
	Priority int    `json:"priority"`
	Status   string `json:"status"`
}

type NodeAllocationState struct {
	NodeId        string   `json:"node_id"`
	AllocationId  string   `json:"allocation_id"`
	LogicId       string   `json:"logic_id"`
	LogicVersion string `json:"logic_version"`
	DesiredStatus string   `json:"desired_status"`
	ActualStatus  string   `json:"actual_status"`
	Trace         []*Trace `json:"trace"`
}

type Trace struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}
