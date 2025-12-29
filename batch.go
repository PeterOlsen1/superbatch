package superbatch

import (
	"fmt"
	"time"

	sp "github.com/PeterOlsen1/superpool"
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
	for _, e := range batchCopy {
		err := b.onFlush(e)
		if err != nil {
			return err
		}
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

	// if threaded, the pool is guaranteed to be created
	if b.threaded {
		return b.flushThreadedUnsafe()
	}
	return b.flushUnsafe()
}

func (b *Batch[T]) flushThreadedUnsafe() error {
	batchCopy := b.copy()

	// TODO: update to batch add later
	for _, e := range batchCopy {
		b.pool.Add(e)
	}
	b.pool.Wait()

	// immediatley return first error, maybe not smart
	for err := range b.pool.Errors() {
		return err
	}

	b.lastFlushed = time.Now()
	b.batch = make([]T, 0, b.cap)
	return nil
}

// Threaded flush.
//
// All items are removed from the batch
// and processed with a workerpool.
func (b *Batch[T]) FlushThreaded() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.pool == nil {
		p, err := sp.NewPool(b.cap, 10, sp.EventHandler[T](b.onFlush))
		if err != nil {
			b.threaded = false
		} else {
			b.pool = p
		}
	}

	return b.flushThreadedUnsafe()
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
	for _, e := range batchCopy {
		err := b.onFlush(e)
		if err != nil {
			return err
		}
	}

	b.lastFlushed = time.Now()
	b.batch = make([]T, 0, b.cap)
	return nil
}
