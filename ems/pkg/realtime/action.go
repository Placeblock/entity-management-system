package realtime

type Action struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
