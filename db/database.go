package db

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
)

const (
	FarFutureEpoch uint64 = 0x7fffffffffffffff
)

type Database struct {
	validators         []*beacon.ValidatorStatus
	validatorPubkeyMap map[beacon.ValidatorPubkey]*beacon.ValidatorStatus

	// Internal fields
	logger *slog.Logger
}

func NewDatabase(logger *slog.Logger) *Database {
	return &Database{
		logger:             logger,
		validators:         []*beacon.ValidatorStatus{},
		validatorPubkeyMap: make(map[beacon.ValidatorPubkey]*beacon.ValidatorStatus),
	}
}

func (db *Database) AddValidator(pubkey beacon.ValidatorPubkey, withdrawalCredentials common.Hash) (*beacon.ValidatorStatus, error) {
	if _, exists := db.validatorPubkeyMap[pubkey]; exists {
		return nil, fmt.Errorf("validator with pubkey %s already exists", pubkey.HexWithPrefix())
	}

	index := len(db.validators)
	validator := &beacon.ValidatorStatus{
		Pubkey:                     pubkey,
		Index:                      strconv.FormatInt(int64(index), 10),
		WithdrawalCredentials:      withdrawalCredentials,
		Balance:                    32e9,
		Status:                     beacon.ValidatorState_PendingInitialized,
		EffectiveBalance:           32e9,
		Slashed:                    false,
		ActivationEligibilityEpoch: 0,
		ActivationEpoch:            0,
		ExitEpoch:                  FarFutureEpoch,
		WithdrawableEpoch:          FarFutureEpoch,
		Exists:                     true,
	}
	db.validators = append(db.validators, validator)
	db.validatorPubkeyMap[pubkey] = validator
	return validator, nil
}

// Get a validator by its index. Returns nil if it doesn't exist.
func (db *Database) GetValidatorByIndex(index uint) *beacon.ValidatorStatus {
	dbLength := len(db.validators)
	if index >= uint(dbLength) {
		return nil
	}

	return db.validators[index]
}

// Get a validator by its pubkey. Returns nil if it doesn't exist.
func (db *Database) GetValidatorByPubkey(pubkey beacon.ValidatorPubkey) *beacon.ValidatorStatus {
	return db.validatorPubkeyMap[pubkey]
}

func (db *Database) Clone() *Database {
	clone := NewDatabase(db.logger)
	cloneValidators := make([]*beacon.ValidatorStatus, len(db.validators))
	for i, validator := range db.validators {
		cloneValidator := &beacon.ValidatorStatus{
			Pubkey:                     validator.Pubkey,
			Index:                      validator.Index,
			WithdrawalCredentials:      validator.WithdrawalCredentials,
			Balance:                    validator.Balance,
			Status:                     validator.Status,
			EffectiveBalance:           validator.EffectiveBalance,
			Slashed:                    validator.Slashed,
			ActivationEligibilityEpoch: validator.ActivationEligibilityEpoch,
			ActivationEpoch:            validator.ActivationEpoch,
			ExitEpoch:                  validator.ExitEpoch,
			WithdrawableEpoch:          validator.WithdrawableEpoch,
			Exists:                     validator.Exists,
		}
		cloneValidators[i] = cloneValidator
		clone.validatorPubkeyMap[validator.Pubkey] = cloneValidator
	}
	clone.validators = cloneValidators
	return clone
}
