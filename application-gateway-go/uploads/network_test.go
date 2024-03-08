/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func AssertNewTestNetwork(t *testing.T, networkName string, options ...ConnectOption) *Network {
	gateway := AssertNewTestGateway(t, options...)
	return gateway.GetNetwork(networkName)
}

func TestNetwork(t *testing.T) {
	t.Run("GetContract returns correctly named Contract", func(t *testing.T) {
		chaincodeName := "basic"
		mockClient := NewMockGatewayClient(gomock.NewController(t))
		network := AssertNewTestNetwork(t, "network", WithGatewayClient(mockClient))

		contract := network.GetContract(chaincodeName)

		require.NotNil(t, contract)
		require.Equal(t, chaincodeName, contract.ChaincodeName(), "chaincode name")
		require.Equal(t, "", contract.ContractName(), "contract name")
	})

	t.Run("GetContractWithName returns correctly named Contract", func(t *testing.T) {
		chaincodeName := "basic"
		contractName := "SimpleChaincode"
		mockClient := NewMockGatewayClient(gomock.NewController(t))
		network := AssertNewTestNetwork(t, "network", WithGatewayClient(mockClient))

		contract := network.GetContractWithName(chaincodeName, contractName)

		require.NotNil(t, contract)
		require.Equal(t, chaincodeName, contract.ChaincodeName(), "basic")
		require.Equal(t, contractName, contract.ContractName(), "SimpleChaincode")
	})
}
