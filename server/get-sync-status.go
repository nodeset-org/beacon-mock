package server

import (
	"net/http"

	"github.com/rocket-pool/node-manager-core/beacon/client"
)

// Handle a get sync status request
func (s *BeaconMockServer) getSyncStatus(w http.ResponseWriter, r *http.Request) {
	// Get the request vars
	_ = s.processApiRequest(w, r, nil)

	// Get the slots
	currentSlot := s.manager.GetCurrentSlot()
	highestSlot := s.manager.GetHighestSlot()

	// Write the response
	response := client.SyncStatusResponse{}
	response.Data.IsSyncing = (currentSlot < highestSlot)
	response.Data.HeadSlot = client.Uinteger(highestSlot)
	response.Data.SyncDistance = client.Uinteger(highestSlot - currentSlot)
	handleSuccess(w, s.logger, response)
}
