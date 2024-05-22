package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/node-manager-core/beacon"
)

type Validator struct {
	Index                 uint64
	Pubkey                beacon.ValidatorPubkey
	WithdrawalCredentials common.Hash
}
