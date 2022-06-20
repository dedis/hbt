package serde

// Serde abstracts away the serialization/deserialization mechanism
type Serde interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(buff []byte, v any) error
}
