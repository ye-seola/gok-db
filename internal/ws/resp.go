package ws

import "encoding/json"

type Response struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
