package registry

// RegistrationData contains the data for a registration
type RegistrationData struct {
	Name       string `json:"name"`
	Passport   string `json:"passport"`
	Picture    []byte `json:"picture"`
	Role       uint64 `json:"role"`
	Registered bool   `json:"registered"`
}

// EncryptedData contains the above encrypted data for a registration
// and a flag to indicate if the data has been successfully registered
type EncryptedData struct {
	Name       []byte `json:"name"`
	Passport   []byte `json:"passport"`
	Picture    []byte `json:"picture"`
	Role       []byte `json:"role"`
	Registered []byte `json:"registered"`
}

// RegistrationID contains the reference to the document in the database
type RegistrationID struct {
	ID []byte `json:"doc_id"`
}
