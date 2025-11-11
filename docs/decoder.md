## Decoder Architecture

### Overview

The Decoder is a pluggable component responsible for transforming raw blockchain event logs into structured, typed events. It operates independently of the Processor, allowing users to either use the provided StandardDecoder implementation or implement custom decoding logic.

### Core Responsibilities

- Parse ABI (Application Binary Interface) definitions
- Match raw logs to event definitions via topic hash lookup
- Decode indexed parameters from log topics
- Decode non-indexed parameters from log data field
- Produce structured Event objects with typed field access

### Interface Definition

The Decoder interface defines the minimal contract for log decoding:

```go
type Decoder interface {
    Decode(log types.Log) (*Event, error)
    DecodeBatch(logs []types.Log) ([]*Event, error)
    GetTopics() []string
}
```

**Methods:**
- `Decode`: Transforms a single raw log into a decoded Event
- `DecodeBatch`: Processes multiple logs efficiently (default implementation may iterate Decode)
- `GetTopics`: Returns registered event signatures for RPC filtering coordination

### Event Structure

Events represent the decoded output with standardized metadata and flexible field storage:

```go
type Event struct {
    // Blockchain Metadata
    ChainID      string
    BlockNumber  uint64
    BlockHash    string
    TxHash       string
    TxIndex      uint64
    LogIndex     uint64
    
    // Contract Information
    Address      string
    
    // Event Identification
    EventType    string   // "Transfer", "Approval", etc.
    EventSig     string   // "Transfer(address,address,uint256)"
    TopicHash    string   // Keccak256 hash of signature
    
    // Decoded Fields
    Fields       EventFields
    
    // Raw Log Reference
    Raw          *types.Log
}

type EventFields map[string]interface{}
```

**EventFields** provides typed accessor methods:
- `GetAddress(key string) (string, error)`
- `GetUint64(key string) (uint64, error)`
- `GetBigInt(key string) (*big.Int, error)`
- `GetBool(key string) (bool, error)`
- `GetBytes(key string) ([]byte, error)`
- `GetString(key string) (string, error)`

### StandardDecoder Implementation

StandardDecoder is the default ABI-based implementation provided by the SDK.

**Internal Structure:**
```go
type StandardDecoder struct {
    events    map[string]*EventDefinition  // topic hash → event definition
    contracts map[string]*ContractABI      // contract address → ABI (optional)
}
```

**Configuration Methods:**
- `RegisterABI(abiJSON string) error` - Parse and register all events from ABI JSON
- `RegisterABIFromFile(filepath string) error` - Load ABI from file
- `AddERC20() error` - Register built-in ERC20 ABI
- `AddERC721() error` - Register built-in ERC721 ABI

### ABI Requirements

The decoder requires full event definitions, not just signatures. Each event definition must include:

1. **Event name** - Human-readable identifier
2. **Parameter names** - Used as keys in Event.Fields map
3. **Parameter types** - Determines decoding logic (address, uint256, etc.)
4. **Indexed status** - Identifies which parameters are in topics vs data field

**Example ABI Event Definition:**
```json
{
  "name": "Transfer",
  "type": "event",
  "inputs": [
    {"name": "from", "type": "address", "indexed": true},
    {"name": "to", "type": "address", "indexed": true},
    {"name": "value", "type": "uint256", "indexed": false}
  ]
}
```

### Decoding Process

1. **Topic Hash Lookup**: Extract `topics[0]` and lookup event definition in registry
2. **Indexed Parameter Decoding**: Decode `topics[1..n]` based on indexed parameter types
3. **Data Field Decoding**: Parse log data field for non-indexed parameters
4. **Field Population**: Combine decoded values into EventFields map
5. **Event Construction**: Return structured Event with metadata and fields

### Type Mapping (Solidity to Go)

| Solidity Type | Go Type | Storage |
|---------------|---------|---------|
| address | string | topics or data |
| uint8-uint64 | uint64 | topics or data |
| uint256 | *big.Int | topics or data |
| int8-int64 | int64 | topics or data |
| int256 | *big.Int | topics or data |
| bool | bool | topics or data |
| bytes, bytesN | []byte | data only |
| string | string | data only |
| arrays | []T | data only |

### Multi-Chain Support

StandardDecoder is chain-agnostic and can be shared across multiple EVM chains:

```go
decoder := core.NewStandardDecoder()
decoder.RegisterABI(erc20ABI)  // Works on all EVM chains

processor.AddChain(ethereumChain, &Options{Decoder: decoder})
processor.AddChain(polygonChain, &Options{Decoder: decoder})
```

**Rationale**: Standard Solidity events use identical encoding across all EVM-compatible chains.

### Multi-Contract Support

A single decoder can handle multiple contracts with different ABIs:

```go
decoder := core.NewStandardDecoder()
decoder.RegisterABI(usdcABI)     // ERC20 events
decoder.RegisterABI(uniswapABI)  // DEX events

// Decoder routes logs to correct ABI via topic hash
```

### Custom Decoder Implementation

Users can implement custom decoders for non-standard requirements:

**Use Cases:**
- Non-EVM chains with different encoding schemes
- Pre-filtering or post-processing logic
- Event enrichment with external data
- Performance optimization for known event structures

**Implementation Requirements:**
- Satisfy Decoder interface (Decode, DecodeBatch, GetTopics)
- Return Event objects with proper structure
- Handle errors gracefully

### Integration with Processor

The Processor invokes the decoder after log fetching:
