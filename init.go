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
		closeChan:     make(chan struct{}),
		batchOpen:     false,
		flushInterval: flushInterval,
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
	b.batchOpen = true

	batchTicker := time.NewTicker(b.flushInterval)
	go func() {
		for {
			select {
			// batch has closed, flush all items and close channels
			case <-b.closeChan:
				b.Flush()
				batchTicker.Stop()
				close(b.closeChan)
				close(b.fullChan)
				return

			// ticker went off or batch is full. flush those items!
			case <-b.fullChan:
			case <-batchTicker.C:
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

	b.closeChan <- struct{}{}
	b.batchOpen = false
	return nil
}
