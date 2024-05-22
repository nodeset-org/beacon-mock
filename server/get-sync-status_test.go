package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/nodeset-org/beacon-mock/api"
	idb "github.com/nodeset-org/beacon-mock/internal/db"
	"github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/stretchr/testify/require"
)

// Make sure sync status requests work
func TestSyncStatus(t *testing.T) {
	currentSlot := uint64(12)

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
	server.manager.SetCurrentSlot(currentSlot)

	// Create the request
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/eth/%s", port, api.SyncingRoute), nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
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

	// Read the body
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("error reading the response body: %v", err)
	}
	var parsedResponse client.SyncStatusResponse
	err = json.Unmarshal(bytes, &parsedResponse)
	if err != nil {
		t.Fatalf("error deserializing response: %v", err)
	}

	// Make sure the response is correct
	require.Equal(t, currentSlot, uint64(parsedResponse.Data.HeadSlot))
	require.Equal(t, uint64(0), uint64(parsedResponse.Data.SyncDistance))
	require.False(t, parsedResponse.Data.IsSyncing)
	t.Logf("Received correct response - head slot: %d, sync distance: %d, is syncing: %t", parsedResponse.Data.HeadSlot, parsedResponse.Data.SyncDistance, parsedResponse.Data.IsSyncing)
}
