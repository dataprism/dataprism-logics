package nodes

type NodeSummary struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type Node struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Resources *NodeResources `json:"resources"`
}

type NodeResources struct {
	CPU *int `json:"cpu"`
	Memory *int `json:"memory_mb"`
	Disk *int `json:"disk_mb"`
	IOPS *int `json:"iops"`
}