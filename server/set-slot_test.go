package server

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/nodeset-org/beacon-mock/api"
	idb "github.com/nodeset-org/beacon-mock/internal/db"
	"github.com/stretchr/testify/require"
)

// Test setting the current slot
func TestSetSlot(t *testing.T) {
	headSlot := uint64(14)
	currentSlot := uint64(9)

	// Take a snapshot
	server.manager.TakeSnapshot("test")
	defer func() {
		err := server.manager.RevertToSnapshot("test")
		if err != nil {
			t.Fatalf("error reverting to snapshot: %v", err)
		}
	}()

	// Provision the database
	d := idb.ProvisionDatabaseForTesting(t, logger)
	server.manager.SetDatabase(d)

	// Send the head slot request
	sendSetSlotRequest(t, headSlot)

	// Send the current slot request
	sendSetSlotRequest(t, currentSlot)

	// Get the sync status now
	parsedResponse := getSyncStatusResponse(t)

	// Make sure the response is correct
	require.Equal(t, headSlot, uint64(parsedResponse.Data.HeadSlot))
	require.Equal(t, (headSlot - currentSlot), uint64(parsedResponse.Data.SyncDistance))
	require.True(t, parsedResponse.Data.IsSyncing)
	t.Logf("Received correct response - head slot: %d, sync distance: %d, is syncing: %t", parsedResponse.Data.HeadSlot, parsedResponse.Data.SyncDistance, parsedResponse.Data.IsSyncing)
}

func sendSetSlotRequest(t *testing.T, slot uint64) {
	// Create the request
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/admin/%s", port, api.SetSlotRoute), nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	query := request.URL.Query()
	query.Add("slot", strconv.FormatUint(slot, 10))
	request.URL.RawQuery = query.Encode()
	t.Logf("Created request")

	// Send the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("error sending request: %v", err)
	}
	t.Logf("Sent request")

	// Check the status code
	require.Equal(t, http.StatusOK, response.StatusCode)
	t.Logf("Received OK status code")
}
