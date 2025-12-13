package superbatch

import (
	"sync"
	"time"
)

type FlushFunc[T any] func([]T) error

type Batch[T any] struct {
	// The underlying sice that holds batch data
	batch []T

	// Mutex for safe locking of slice in goroutines
	mu sync.Mutex

	// Capacity of the batch, need to keep track for when batch is reset
	cap uint32

	// Signal when the batch is full, will be emptied after
	fullChan chan struct{}

	// Signals to the current running goroutine to stop
	stopChan chan struct{}

	batchOpen bool

	// Function passed in from the init method
	//
	// When flushing, this function will be applied to all
	// members of the batch
	onFlush FlushFunc[T]

	// The time the batch was last flushed
	lastFlushed time.Time

	flushInterval time.Duration

	ticker *time.Ticker
}
