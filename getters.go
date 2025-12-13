package superbatch

func (b *Batch[T]) Len() int {
	return len(b.batch)
}
