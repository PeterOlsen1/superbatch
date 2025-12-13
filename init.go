package superbatch

import (
	"fmt"
	"time"
)

func InitBatch[T any](cap uint32, flushInterval time.Duration, onFlush FlushFunc[T]) (*Batch[T], error) {
	if cap == 0 {
		return nil, fmt.Errorf("capacity cannot be 0")
	}

	b := &Batch[T]{
		batch:         make([]T, 0, cap),
		cap:           cap,
		onFlush:       onFlush,
		fullChan:      make(chan struct{}),
		batchOpen:     false,
		flushInterval: flushInterval,
		ticker:        nil,
	}

	b.Open()
	return b, nil
}

// Open the batch
//
// This function initializes the batch ticker
// and channel listener for the close signal
//
// Returns error when called on open batch
func (b *Batch[T]) Open() error {
	if b.batchOpen {
		return fmt.Errorf("batch is already open")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	err := b.startTicker()
	if err != nil {
		return err
	}

	b.batchOpen = true
	return nil
}

// Start the batch ticker
//
// It is assumed that this is called with the mutex locked
// It is also assumed that this is not called more than once
func (b *Batch[T]) startTicker() error {
	if b.ticker != nil {
		return fmt.Errorf("ticker is not nil")
	}

	b.ticker = time.NewTicker(b.flushInterval)
	go func() {
		for {
			select {
			// ticker went off or batch is full. flush those items!
			case <-b.fullChan:
			case <-b.ticker.C:
				b.Flush()
			}
		}
	}()

	return nil
}

// Closes the batch.
//
// Once this signal is sent, the ticker will stop
func (b *Batch[T]) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.batchOpen {
		return fmt.Errorf("batch is already closed")
	}

	b.flushUnsafe()
	b.ticker.Stop()
	b.ticker = nil
	close(b.fullChan)
	b.batchOpen = false
	return nil
}
