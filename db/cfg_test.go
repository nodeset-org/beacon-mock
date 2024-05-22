package db

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nodeset-org/beacon-mock/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConfigClone(t *testing.T) {
	depositContract := common.HexToAddress(test.DepositContractAddressString)
	c := NewConfig(depositContract, true)
	clone := c.Clone()
	t.Log("Created config and clone")

	require.NotSame(t, c, clone)
	require.Equal(t, c, clone)
	t.Log("Configs are equal")
}
