package realtime

import "encoding/json"

type RawAction struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
