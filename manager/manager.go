package manager

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/beacon-mock/db"
	"github.com/rocket-pool/node-manager-core/beacon"
)

// Beacon mock manager
type BeaconMockManager struct {
	database *db.Database
	config   *db.Config

	// Internal fields
	snapshots map[string]*db.Database
	logger    *slog.Logger
}

// Create a new beacon mock manager instance
func NewBeaconMockManager(logger *slog.Logger, config *db.Config) *BeaconMockManager {
	return &BeaconMockManager{
		database:  db.NewDatabase(logger),
		config:    config,
		snapshots: map[string]*db.Database{},
		logger:    logger,
	}
}

// Set the database for the manager directly if you need to custom provision it
func (m *BeaconMockManager) SetDatabase(db *db.Database) {
	m.database = db
}

// Take a snapshot of the current database state
func (m *BeaconMockManager) TakeSnapshot(name string) {
	m.snapshots[name] = m.database.Clone()
	m.logger.Info("Took DB snapshot", "name", name)
}

// Revert to a snapshot of the database state
func (m *BeaconMockManager) RevertToSnapshot(name string) error {
	snapshot, exists := m.snapshots[name]
	if !exists {
		return fmt.Errorf("snapshot with name [%s] does not exist", name)
	}
	m.database = snapshot
	m.logger.Info("Reverted to DB snapshot", "name", name)
	return nil
}

// Returns the manager's Beacon config
func (m *BeaconMockManager) GetConfig() *db.Config {
	return m.config
}

// Increments the Beacon chain slot - use this to have parity with mining a new block on the EL
func (m *BeaconMockManager) IncrementSlot() {
	m.database.SetCurrentSlot(m.database.GetCurrentSlot() + 1)
}

// Returns the current Beacon chain slot
func (m *BeaconMockManager) GetCurrentSlot() uint64 {
	return m.database.GetCurrentSlot()
}

// Returns the highest Beacon chain slot (top of the chain head)
func (m *BeaconMockManager) GetHighestSlot() uint64 {
	return m.database.GetHighestSlot()
}

// Sets the current Beacon chain slot
func (m *BeaconMockManager) SetCurrentSlot(slot uint64) {
	m.database.SetCurrentSlot(slot)
}

// Add a validator to the Beacon chain
func (m *BeaconMockManager) AddValidator(pubkey beacon.ValidatorPubkey, withdrawalCredentials common.Hash) (*db.Validator, error) {
	return m.database.AddValidator(pubkey, withdrawalCredentials)
}

// Gets a validator by its index or pubkey
func (m *BeaconMockManager) GetValidator(id string) (*db.Validator, error) {
	if len(id) == beacon.ValidatorPubkeyLength*2 || strings.HasPrefix(id, "0x") {
		pubkey, err := beacon.HexToValidatorPubkey(id)
		if err != nil {
			return nil, fmt.Errorf("error parsing pubkey [%s]: %v", id, err)
		}
		return m.database.GetValidatorByPubkey(pubkey), nil
	}
	index, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing index [%s]: %v", id, err)
	}
	return m.database.GetValidatorByIndex(uint(index)), nil
}

// Gets multiple validators by their indices or pubkeys
func (m *BeaconMockManager) GetValidators(ids []string) ([]*db.Validator, error) {
	if len(ids) == 0 {
		return m.database.GetAllValidators(), nil
	}

	validators := []*db.Validator{}
	for _, id := range ids {
		validator, err := m.GetValidator(id)
		if err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}
	return validators, nil
}
