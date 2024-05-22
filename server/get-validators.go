package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/nodeset-org/beacon-mock/api"
	"github.com/rocket-pool/node-manager-core/beacon/client"
	"github.com/rocket-pool/node-manager-core/log"
)

// Handle a get validators request
func (s *BeaconMockServer) getValidators(w http.ResponseWriter, r *http.Request) {
	// Get the request vars
	args := s.processApiRequest(w, r, nil)

	var ids []string
	switch r.Method {
	case http.MethodGet:
		ids = s.getValidatorIDsFromRequestArgs(args)
	case http.MethodPost:
		ids = s.getValidatorIDsFromRequestBody(w, r)
		if ids == nil {
			return
		}
	default:
		handleInvalidMethod(s.logger, w)
		return
	}

	// Get the validators
	validators, err := s.manager.GetValidators(ids)
	if err != nil {
		handleInputError(s.logger, w, err)
		return
	}

	// Write the response
	validatorMetas := make([]client.Validator, len(validators))
	for i, validator := range validators {
		validatorMetas[i] = getValidatorMetaFromValidator(validator)
	}
	response := client.ValidatorsResponse{
		Data: validatorMetas,
	}
	handleSuccess(w, s.logger, response)
}

// Get all of the validator IDs from the request query for a GET request
func (s *BeaconMockServer) getValidatorIDsFromRequestArgs(args url.Values) []string {
	ids := args["id"]
	return getValidatorIDs(ids)
}

// Get all of the validator IDs from the request body for a POST request
func (s *BeaconMockServer) getValidatorIDsFromRequestBody(w http.ResponseWriter, r *http.Request) []string {
	// Read the body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		handleInputError(s.logger, w, fmt.Errorf("error reading request body: %w", err))
		return nil
	}
	s.logger.Debug("Request body:", slog.String(log.BodyKey, string(bodyBytes)))

	// Deserialize the body
	var requestBody api.ValidatorsRequest
	err = json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		handleInputError(s.logger, w, fmt.Errorf("error deserializing request body: %w", err))
		return nil
	}

	return getValidatorIDs(requestBody.IDs)
}

// Get all of the validator IDs from a list of them, handling the case where they're comma-separated
func getValidatorIDs(ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}

	fullIds := make([]string, 0, len(ids))
	for _, id := range ids {
		elements := strings.Split(id, ",")
		for _, element := range elements {
			trimmed := strings.TrimSpace(element)
			if trimmed == "" {
				continue
			}
			fullIds = append(fullIds, trimmed)
		}
	}
	return fullIds
}
