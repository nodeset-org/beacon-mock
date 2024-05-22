package db

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
)

// Beacon mock database
type Database struct {
	// Validators registered with the network
	validators []*Validator

	// Lookup of validators by pubkey
	validatorPubkeyMap map[beacon.ValidatorPubkey]*Validator

	// Current slot
	currentSlot uint64

	// Highest slot
	highestSlot uint64

	// Internal fields
	logger *slog.Logger
}

// Create a new database instance
func NewDatabase(logger *slog.Logger) *Database {
	return &Database{
		logger:             logger,
		validators:         []*Validator{},
		validatorPubkeyMap: make(map[beacon.ValidatorPubkey]*Validator),
	}
}

// Add a new validator to the database. Returns an error if the validator already exists.
func (db *Database) AddValidator(pubkey beacon.ValidatorPubkey, withdrawalCredentials common.Hash) (*Validator, error) {
	if _, exists := db.validatorPubkeyMap[pubkey]; exists {
		return nil, fmt.Errorf("validator with pubkey %s already exists", pubkey.HexWithPrefix())
	}

	index := len(db.validators)
	validator := NewValidator(pubkey, withdrawalCredentials, uint64(index))
	db.validators = append(db.validators, validator)
	db.validatorPubkeyMap[pubkey] = validator
	return validator, nil
}

// Get a validator by its index. Returns nil if it doesn't exist.
func (db *Database) GetValidatorByIndex(index uint) *Validator {
	dbLength := len(db.validators)
	if index >= uint(dbLength) {
		return nil
	}

	return db.validators[index]
}

// Get a validator by its pubkey. Returns nil if it doesn't exist.
func (db *Database) GetValidatorByPubkey(pubkey beacon.ValidatorPubkey) *Validator {
	return db.validatorPubkeyMap[pubkey]
}

// Get all validators
func (db *Database) GetAllValidators() []*Validator {
	return db.validators
}

// Get the latest slot
func (db *Database) GetCurrentSlot() uint64 {
	return db.currentSlot
}

// Get the highest slot on the chain
func (db *Database) GetHighestSlot() uint64 {
	return db.highestSlot
}

// Set the current slot - this will also update the highest slot if the new slot is higher
func (db *Database) SetCurrentSlot(slot uint64) {
	db.currentSlot = slot
	if slot > db.highestSlot {
		db.highestSlot = slot
	}
}

// Clone the database into a new instance
func (db *Database) Clone() *Database {
	clone := NewDatabase(db.logger)
	clone.currentSlot = db.currentSlot
	clone.highestSlot = db.highestSlot
	cloneValidators := make([]*Validator, len(db.validators))
	for i, validator := range db.validators {
		cloneValidator := validator.Clone()
		cloneValidators[i] = cloneValidator
		clone.validatorPubkeyMap[validator.Pubkey] = cloneValidator
	}
	clone.validators = cloneValidators
	return clone
}
