package manager

import (
	"fmt"
	"log/slog"

	"github.com/nodeset-org/beacon-mock/db"
)

type BeaconMockManager struct {
	Database *db.Database

	// Internal fields
	snapshots map[string]*db.Database
	logger    *slog.Logger
}

func NewBeaconMockManager(logger *slog.Logger) *BeaconMockManager {
	return &BeaconMockManager{
		Database:  db.NewDatabase(logger),
		snapshots: map[string]*db.Database{},
		logger:    logger,
	}
}

func (m *BeaconMockManager) TakeSnapshot(name string) {
	m.snapshots[name] = m.Database.Clone()
	m.logger.Info("Took DB snapshot", "name", name)
}

func (m *BeaconMockManager) RevertToSnapshot(name string) error {
	snapshot, exists := m.snapshots[name]
	if !exists {
		return fmt.Errorf("snapshot with name [%s] does not exist", name)
	}
	m.Database = snapshot
	m.logger.Info("Reverted to DB snapshot", "name", name)
	return nil
}
