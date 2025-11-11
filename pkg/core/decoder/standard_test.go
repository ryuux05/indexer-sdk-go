package decoder

import (
	//"encoding/json"
	//"fmt"
	"math/big"
	"testing"

	"github.com/ryuux05/godex/pkg/core/types"
	"github.com/stretchr/testify/assert"
)

const erc20Transfer_ABI = `[
	{
	  "anonymous": false,
	  "inputs": [
		{
		  "indexed": true,
		  "internalType": "address",
		  "name": "from",
		  "type": "address"
		},
		{
		  "indexed": true,
		  "internalType": "address",
		  "name": "to",
		  "type": "address"
		},
		{
		  "indexed": false,
		  "internalType": "uint256",
		  "name": "value",
		  "type": "uint256"
		}
	  ],
	  "name": "Transfer",
	  "type": "event"
	}
  ]`

func TestDecodeTransfer_Successful(t *testing.T) {
	decoder := NewStandsardDecoder()
	decoder.RegisterABI(erc20Transfer_ABI)

	log := types.Log{
		Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Topics: []string{
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"0x000000000000000000000000a1b2c3d4e5f6789012345678901234567890abcd",
			"0x000000000000000000000000f1e2d3c4b5a6978012345678901234567890dcba",
		},
		Data:             "0x0000000000000000000000000000000000000000000000000000000005f5e100",
		BlockNumber:      "0x112a880",
		BlockHash:        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		TransactionHash:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TransactionIndex: "0x0",
		LogIndex:         "0x5",
	}

	event, err := decoder.Decode(log)

	// Print as JSON
	//eventsJSON, _ := json.MarshalIndent(decoder.events, "", "  ")
	//fmt.Printf("Registered Events:\n%s\n", string(eventsJSON))

	assert.NoError(t, err, "Error decoding log")
	assert.Equal(t, event.BlockNumber, uint64(18000000))
	assert.Equal(t, event.Address, "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	assert.Equal(t, event.BlockHash, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	assert.Equal(t, event.TransactionHash, "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	assert.Equal(t, event.LogIndex, uint64(5))
	assert.Equal(t, event.EventType, "Transfer")
	assert.Equal(t, event.Fields["from"], "0xa1b2c3d4e5f6789012345678901234567890abcd")
	assert.Equal(t, event.Fields["to"], "0xf1e2d3c4b5a6978012345678901234567890dcba")
	assert.Equal(t, event.Fields["value"], big.NewInt(100000000))

}
