package logics

type Logic struct {
	Id string `json:"id"`
	Description string `json:"description"`
	Language string `json:"language"`
	Code string `json:"code"`
	Libraries []string `json:"libraries"`
	Resources *LogicResources `json:"resources"`
	InboundTopics []string `json:"inbound_topics"`
	OutboundTopics []string `json:"outbound_topics"`
}

type LogicResources struct {
	CPU *int `json:"cpu_mhz"`
	Memory *int `json:"memory_mb"`
	Disk *int `json:"disk_mb"`
}

type LogicStatus struct {
	Queued int `json:"queued"`
	Complete int `json:"complete"`
	Failed int `json:"failed"`
	Running int `json:"running"`
	Starting int `json:"starting"`
	Lost int `json:"lost"`
}
