# indexer-sdk-go

A high-performance blockchain indexing SDK written in Go, designed for building robust indexers that process EVM-compatible blockchain events with automatic reorg handling, multi-chain support, and flexible event decoding.

## Features

- **Multi-Chain Support**: Index events across multiple EVM-compatible chains simultaneously
- **Automatic Reorg Handling**: Built-in detection and rollback for blockchain reorganizations
- **Flexible Event Decoding**: Register multiple ABIs with named identifiers for explicit event decoding
- **High Performance**: Concurrent fetching and processing with configurable worker pools
- **Production Ready**: Comprehensive error handling, retry mechanisms, and state management

## Installation

```bash
go get github.com/ryuux05/godex
```

## Quick Start

### Basic Indexer Setup

```go
package main

import (
    "context"
    "log"
    
    "github.com/ryuux05/godex/pkg/core"
    "github.com/ryuux05/godex/pkg/core/decoder"
)

func main() {
    // Initialize RPC client
    rpc := core.NewHTTPRPC("https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY", 20)
    
    // Configure processor options
    opts := core.Options{
        RangeSize:          100,
        BatchSize:          50,
        DecoderConcurrency: 4,
        FetcherConcurrency: 4,
        StartBlock:         18000000,
        Confimation:        15,
        LogsBufferSize:     1024,
        Topics: []string{
            "Transfer(address,address,uint256)",
        },
        FetchMode: core.FetchModeReceipts,
    }
    
    // Setup chain
    chain := core.ChainInfo{
        ChainId: "1",
        Name:    "Ethereum",
        RPC:     rpc,
    }
    
    // Initialize decoder and register ABIs
    decoder := decoder.NewStandardDecoder()
    decoder.RegisterABI("ERC20", erc20ABI)
    
    // Create processor
    processor := core.NewProcessor()
    processor.AddChain(chain, &opts)
    
    // Start indexing
    ctx := context.Background()
    if err := processor.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Event Decoding

```go
// Register multiple ABIs with identifiers
decoder := decoder.NewStandardDecoder()
decoder.RegisterABI("ERC20", erc20TransferABI)
decoder.RegisterABI("ERC721", erc721TransferABI)

// Get logs from processor
logsCh, err := processor.Logs(chain.ChainId)
if err != nil {
    log.Fatal(err)
}

// Process logs
for log := range logsCh {
    // Select ABI based on log structure
    var event *types.Event
    if len(log.Topics) == 3 {
        event, err = decoder.Decode("ERC20", log)
    } else if len(log.Topics) == 4 {
        event, err = decoder.Decode("ERC721", log)
    }
    
    if err != nil {
        log.Printf("Error decoding: %v", err)
        continue
    }
    
    if event != nil {
        // Access decoded event fields
        from := event.Fields["from"].(string)
        to := event.Fields["to"].(string)
        value := event.Fields["value"].(*big.Int)
        
        log.Printf("Transfer: %s -> %s, value: %s", from, to, value.String())
    }
}
```

### Multi-Chain Indexing

```go
processor := core.NewProcessor()

// Add Ethereum chain
ethereumChain := core.ChainInfo{
    ChainId: "1",
    Name:    "Ethereum",
    RPC:     ethereumRPC,
}
processor.AddChain(ethereumChain, &ethereumOpts)

// Add Polygon chain
polygonChain := core.ChainInfo{
    ChainId: "137",
    Name:    "Polygon",
    RPC:     polygonRPC,
}
processor.AddChain(polygonChain, &polygonOpts)

// Run all chains concurrently
processor.Run(ctx)
```

## Configuration

### Processor Options

- `RangeSize`: Number of blocks to fetch per batch
- `BatchSize`: Number of logs to process per batch
- `DecoderConcurrency`: Number of concurrent decoder workers
- `FetcherConcurrency`: Number of concurrent RPC fetchers
- `StartBlock`: Initial block number to start indexing
- `Confimation`: Number of confirmations required before processing
- `LogsBufferSize`: Buffer size for log channel
- `Topics`: Event signatures to filter (supports function signatures or topic hashes)
- `FetchMode`: Log fetching strategy (`FetchModeLogs` or `FetchModeReceipts`)

### RPC Configuration

```go
rpc := core.NewHTTPRPC(
    "https://your-rpc-endpoint.com",
    20, // Rate limit (requests per second)
)
```

## Architecture

The SDK consists of three main components:

### Processor

The Processor manages block fetching, log retrieval, and reorg detection. It coordinates multiple workers to fetch logs concurrently while maintaining ordered processing through an arbiter pattern.

### Decoder

The Decoder transforms raw blockchain logs into structured events. It supports registering multiple ABIs with named identifiers, allowing explicit selection of which ABI to use for decoding.

### RPC Client

The RPC client handles communication with blockchain nodes, including automatic retry logic, rate limiting, and error handling.

## Reorg Handling

The SDK automatically detects blockchain reorganizations by comparing parent block hashes. When a reorg is detected:

1. The processor identifies the common ancestor block
2. Rolls back any processed state beyond the ancestor
3. Restarts indexing from the ancestor block

This ensures data consistency even during chain reorganizations.

## Error Handling

The decoder is designed to be resilient. It returns `nil, nil` for logs that cannot be decoded (structure mismatches, missing data, etc.), allowing the indexer to continue processing. Only configuration errors (such as ABI not found) return actual errors.

## Performance Considerations

- Use appropriate `FetcherConcurrency` and `DecoderConcurrency` values based on your RPC rate limits
- Adjust `RangeSize` to balance between RPC call frequency and memory usage
- Use `FetchModeReceipts` for better performance when filtering by contract addresses
- Monitor log channel buffer size to prevent blocking

