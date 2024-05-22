package api

const (
	StateID     string = "state_id"
	ValidatorID string = "validator_id"

	// Beacon API routes
	ValidatorsRoute string = "v1/beacon/states/{" + StateID + "}/validators"
	ValidatorRoute  string = "v1/beacon/states/{" + StateID + "}/validators/{" + ValidatorID + "}"
	SyncingRoute    string = "v1/node/syncing"

	// Admin routes
	AddValidatorRoute string = "add-validator"
	SetBalanceRoute   string = "set-balance"
	SetStatusRoute    string = "set-status"
	SetSlotRoute      string = "set-slot"
	SlashRoute        string = "slash"
)
