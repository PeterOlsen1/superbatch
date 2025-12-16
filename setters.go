package superbatch

import "time"

// Sets the current flush interval.
//
// If the batch is open: flush, stop ticker, update
func (b *Batch[T]) SetFlushInterval(newInterval *time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// quick return, set flush interval and DON'T open the batch.
	// let the user do that
	if !b.batchOpen {
		b.flushInterval = newInterval
		return nil
	}

	if newInterval == nil {
		b.stopChan <- struct{}{}
		b.ticker = nil
		b.flushInterval = newInterval
		return b.startTicker()
	}

	// new interval is shorter than the time we are at since last flush
	if b.lastFlushed.Add(*newInterval).Before(time.Now()) {
		err := b.flushUnsafe()
		if err != nil {
			return err
		}
		b.lastFlushed = time.Now()
	}

	b.stopChan <- struct{}{}
	b.ticker = nil
	b.flushInterval = newInterval
	return b.startTicker() // automatically handles nil newInterval case
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

	b.cap = newCap
	// we have enough space, just copy over the items
	if uint32(len(b.batch)) < newCap {
		itemsCopy := b.copy()
		b.batch = make([]T, 0, newCap)
		b.batch = append(b.batch, itemsCopy...)
		return nil
	}

	err := b.flushUnsafe()
	if err != nil {
		return err
	}
	b.batch = make([]T, 0, newCap)
	return nil
}
