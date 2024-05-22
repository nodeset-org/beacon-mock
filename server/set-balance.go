package server

import "github.com/rocket-pool/node-manager-core/beacon"

func (s *BeaconMockServer) SetBalance(validator *beacon.ValidatorStatus, newBalanceGwei uint64) error {
	validator.Balance = newBalanceGwei
	return nil
}
