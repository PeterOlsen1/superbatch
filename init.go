package superbatch

import (
	"fmt"
	"time"

	sp "github.com/PeterOlsen1/superpool"
)

func NewBatch[T any](cfg BatchConfig[T]) (*Batch[T], error) {
	if cfg.Cap == 0 {
		return nil, fmt.Errorf("capacity cannot be 0")
	}

	b := &Batch[T]{
		batch:         make([]T, 0, cfg.Cap),
		cap:           cfg.Cap,
		onFlush:       cfg.OnFlush,
		fullChan:      make(chan struct{}),
		stopChan:      make(chan struct{}),
		batchOpen:     false,
		flushInterval: cfg.FlushInterval,
		ticker:        nil,
		threaded:      cfg.Threaded,
	}

	if cfg.Threaded {
		p, err := sp.NewPool(cfg.Cap, 10, sp.EventHandler[T](cfg.OnFlush))
		if err != nil {
			b.threaded = false
		} else {
			b.pool = p
		}
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
// # The batch ticker is automatically set once started
//
// It is assumed that this is called with the mutex locked
// It is also assumed that this is not called more than once
func (b *Batch[T]) startTicker() error {
	if b.ticker != nil {
		return fmt.Errorf("ticker is not nil")
	}

	if b.flushInterval != nil {
		b.ticker = time.NewTicker(*b.flushInterval)
		go func() {
			for {
				select {
				case <-b.stopChan:
					return
				// ticker went off or batch is full. flush those items!
				case <-b.fullChan:
				case <-b.ticker.C:
					b.Flush()
				}
			}
		}()
	} else {
		go func() {
			for {
				select {
				case <-b.stopChan:
					return
				case <-b.fullChan:
					b.Flush()
				}
			}
		}()
	}

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

	b.stopChan <- struct{}{}
	b.flushUnsafe()

	if b.ticker != nil {
		b.ticker.Stop()
	}
	b.ticker = nil

	if b.pool != nil {
		b.pool.Shutdown()
	}

	close(b.fullChan)
	close(b.stopChan)
	b.batchOpen = false
	return nil
}
