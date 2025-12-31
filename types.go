package superbatch

import (
	"sync"
	"time"

	sp "github.com/PeterOlsen1/superpool"
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

	// Dictates how often the batch is flushed.
	// Pass in nil for no intervals
	flushInterval *time.Duration

	ticker *time.Ticker

	// Config variable to be set if batches should be processed with worker pools
	threaded bool

	// pool to process batch with multiple threads if requested
	pool *sp.Pool[[]T]

	// chan that the user can subscribe to, signaling when last flush happened
	flushed chan struct{}
}

type BatchConfig[T any] struct {
	Cap           uint32
	FlushInterval *time.Duration
	OnFlush       FlushFunc[T]
	Threaded      bool
}
