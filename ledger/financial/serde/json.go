package serde

import "encoding/json"

// JSON implements a JSON-based serialization/deserialization mechanism
//
// - implements state.Serde
type JSON struct{}

// Marshal implements state.Serde
func (JSON) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal implements state.Serde
func (JSON) Unmarshal(buff []byte, v any) error {
	return json.Unmarshal(buff, v)
}
