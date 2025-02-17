package realtime

type RawAction struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
