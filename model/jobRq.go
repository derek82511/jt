package model

type JobRq struct {
	ScenarioID  string `json:"scenarioId"`
	MinHeap     string `json:"minHeap"`
	MaxHeap     string `json:"maxHeap"`
	Config      string `json:"config"`
	ExecuteType int    `json:"executeType"`
	RemoteHost  string `json:"remoteHost"`
}
