package superbatch

import (
	"fmt"
	"time"
)

// Copies the contents of a batch array
func (b *Batch[T]) copy() []T {
	return append([]T(nil), b.batch...)
}

// Add an item to the given batch
//
// If the batch hits capacity, it will be flushed
func (b *Batch[T]) Add(item T) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.batchOpen {
		return fmt.Errorf("batch is closed")
	}

	if len(b.batch) == cap(b.batch)-1 {
		b.batch = append(b.batch, item)
		return b.flushUnsafe()
	}

	b.batch = append(b.batch, item)
	return nil
}

// Unsafe flush.
//
// This means all items will be removed from the batch
// and the mutex will NOT be locked to do so.
// Therefore, this is for internal use only.
//
// The batch will be reset after the operation if successful.
//
// If failed, the batch will NOT be reset, and an error returned.
func (b *Batch[T]) flushUnsafe() error {
	batchCopy := b.copy()
	err := b.onFlush(batchCopy)
	if err != nil {
		return err
	}

	b.lastFlushed = time.Now()
	b.batch = make([]T, 0, b.cap)
	return nil
}

// Safe flush.
//
// All items are removed from the batch
// and the mutex WILL be locked
func (b *Batch[T]) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.flushUnsafe()
}

// Flush all items according to custom function
//
// Custom function still needs to follow the same conventions
// as the regular batch flush function
func (b *Batch[T]) FlushCustom(onFlush FlushFunc[T]) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	//apply custom flush func
	batchCopy := b.copy()
	err := onFlush(batchCopy)
	if err != nil {
		return err
	}

	b.lastFlushed = time.Now()
	b.batch = make([]T, 0, b.cap)
	return nil
}
