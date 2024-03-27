package smc

import "go.dedis.ch/hbt/server/registry/registry"

// SmcSecret contains the secret for the SMC
type Secret struct {
	Data string                  `json:"secret"`
	ID   registry.RegistrationID `json:"id"`
}
