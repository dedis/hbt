package smc

import "go.dedis.ch/hbt/server/registration/registry"

// SmcSecret contains the secret for the SMC
type Secret struct {
	Data string                  `json:"secret"`
	Id   registry.RegistrationID `json:"id"`
}
