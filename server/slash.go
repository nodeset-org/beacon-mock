package server

import (
	"fmt"

	"github.com/rocket-pool/node-manager-core/beacon"
)

func (s *BeaconMockServer) SlashValidator(validator *beacon.ValidatorStatus, penaltyGwei uint64) error {
	if validator.Status != beacon.ValidatorState_ActiveOngoing && validator.Status != beacon.ValidatorState_ActiveExiting {
		return fmt.Errorf("validator with pubkey %s is not in a slashable state", validator.Pubkey.HexWithPrefix())
	}
	validator.Slashed = true
	validator.Balance -= penaltyGwei
	validator.Status = beacon.ValidatorState_ActiveSlashed
	return nil
}
