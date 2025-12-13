package superbatch

import "time"

// Sets the current flush interval.
//
// If the batch is open: flush, stop ticker, update
func (b *Batch[T]) SetFlushInterval(newInterval time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// quick return, set flush interval and DON'T open the batch.
	// let the user do that
	if !b.batchOpen {
		b.flushInterval = newInterval
		return nil
	}

	err := b.flushUnsafe()
	if err != nil {
		return err
	}

	b.stopChan <- struct{}{}
	b.ticker = nil
	b.flushInterval = newInterval
	b.startTicker() // won't error, ticker is set to nil

	return nil
}

// Sets the capacity for the batch.
//
// If the current size of the batch is less than the new capacity,
// all current items will be flushed.
//
// An error is returned if this flush errors
func (b *Batch[T]) SetCap(newCap uint32) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// quick return, set flush interval and DON'T open the batch.
	// let the user do that
	if !b.batchOpen {
		b.cap = newCap
		return nil
	}

	// we have enough space, just copy over the items
	if uint32(len(b.batch)) < newCap {
		items := b.copy()
		b.batch = make([]T, 0, newCap)
		b.batch = append(b.batch, items...)
		return nil
	}

	err := b.flushUnsafe()
	if err != nil {
		return err
	}
	b.batch = make([]T, 0, newCap)
	return nil
}
