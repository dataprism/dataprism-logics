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
