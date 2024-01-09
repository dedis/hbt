package registry

// RegistrationData contains the data for a registration
type RegistrationData struct {
	Name       string `json:"name"`
	Passport   string `json:"passport"`
	Role       uint   `json:"role"`
	Picture    []byte `json:"picture"`
	Hash       []byte `json:"hash"`
	Registered bool   `json:"registered"`
}

type DocID []byte

// RegistrationID contains the reference to the document in the database
type RegistrationID struct {
	ID DocID `json:"doc_id"`
}
