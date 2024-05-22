package server

import "github.com/rocket-pool/node-manager-core/beacon"

func (s *BeaconMockServer) SetStatus(validator *beacon.ValidatorStatus, newStatus beacon.ValidatorState) error {
	validator.Status = newStatus
	return nil
}
