package superbatch

import "time"

func (b *Batch[T]) Len() int {
	return len(b.batch)
}

func (b *Batch[T]) GetLastFlushed() time.Time {
	return b.lastFlushed
}
