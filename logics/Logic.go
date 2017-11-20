package logics

type Logic struct {
	Id string `json:"id"`
	Description string `json:"description"`
}

type LogicVersion struct {
	Version int `json:"version"`
	Language string `json:"language"`
	Code string `json:"code"`
}

type LogicStatus struct {
	Queued int `json:"queued"`
	Complete int `json:"complete"`
	Failed int `json:"failed"`
	Running int `json:"running"`
	Starting int `json:"starting"`
	Lost int `json:"lost"`
}
