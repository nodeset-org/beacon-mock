package server

import (
	"strconv"

	"github.com/nodeset-org/beacon-mock/db"
	"github.com/rocket-pool/node-manager-core/beacon/client"
)

func getValidatorMetaFromValidator(validator *db.Validator) client.Validator {
	validatorMeta := client.Validator{
		Index:   strconv.FormatUint(validator.Index, 10),
		Balance: client.Uinteger(validator.Balance),
		Status:  string(validator.Status),
	}
	validatorMeta.Validator.Pubkey = validator.Pubkey[:]
	validatorMeta.Validator.WithdrawalCredentials = validator.WithdrawalCredentials[:]
	validatorMeta.Validator.EffectiveBalance = client.Uinteger(validator.EffectiveBalance)
	validatorMeta.Validator.Slashed = validator.Slashed
	validatorMeta.Validator.ActivationEligibilityEpoch = client.Uinteger(validator.ActivationEligibilityEpoch)
	validatorMeta.Validator.ActivationEpoch = client.Uinteger(validator.ActivationEpoch)
	validatorMeta.Validator.ExitEpoch = client.Uinteger(validator.ExitEpoch)
	validatorMeta.Validator.WithdrawableEpoch = client.Uinteger(validator.WithdrawableEpoch)
	return validatorMeta
}
