package sink

import "context"

type Sink interface {
	// Store persists events from a single block
	Store(ctx context.Context, chainId string)
	// StoreBatch persists events from multiple blocks efficiently
	// This is useful for batch operations and better performance
	StoreBatch()
	// Rollback removes all events from a block number onwards
	// Used during reorg handling to remove orphaned blocks
	Rollback()
	// GetLastBlock returns the last processed block number for a chain
	// Returns 0 if no blocks have been processed yet
	GetLastBlock()
}
