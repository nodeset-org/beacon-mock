package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/beacon-mock/manager"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type BeaconMockServer struct {
	logger  *slog.Logger
	ip      string
	port    uint16
	socket  net.Listener
	server  http.Server
	router  *mux.Router
	manager *manager.BeaconMockManager
}

// =============
// === Utils ===
// =============

func (s *BeaconMockServer) getValidatorByID(id string) (*beacon.ValidatorStatus, error) {
	if strings.HasPrefix(id, "0x") {
		pubkey, err := beacon.HexToValidatorPubkey(id)
		if err != nil {
			return nil, err
		}
		validator := s.manager.Database.GetValidatorByPubkey(pubkey)
		if validator == nil {
			return nil, fmt.Errorf("validator with pubkey %s does not exist", pubkey.HexWithPrefix())
		}
		return validator, nil
	}

	index, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	validator := s.manager.Database.GetValidatorByIndex(uint(index))
	if validator == nil {
		return nil, fmt.Errorf("validator with index %d does not exist", index)
	}
	return validator, nil
}
