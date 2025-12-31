package superbatch

import "time"

func (b *Batch[T]) Len() int {
	return len(b.batch)
}

func (b *Batch[T]) GetLastFlushed() time.Time {
	return b.lastFlushed
}

func (b *Batch[T]) Flushed() <-chan struct{} {
	if b.flushed == nil {
		b.flushed = make(chan struct{}, 1000)
	}

	return b.flushed
}
