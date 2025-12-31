package superbatch_test

import (
	"testing"
	"time"

	sb "github.com/PeterOlsen1/superbatch"
)

// Initialize a batch of size 10 with no timeout.
// If creation fails, test fails.
func setup(t *testing.T) *sb.Batch[int] {
	count = 0
	batchFunc := func(items []int) error {
		count += 1
		return nil
	}

	cfg := sb.BatchConfig[int]{
		Cap:     10,
		OnFlush: batchFunc,
	}
	b, err := sb.NewBatch(cfg)
	if err != nil {
		t.Fatalf("Failed to setup batch: %s", err)
		return nil
	}
	return b
}

// Initialize a batch with size 10 and timeout
//
// The interval is not accepted as a pointer, since the null case is
// already handled in the default setup() method.
func setupWithInterval(t *testing.T, i time.Duration) *sb.Batch[int] {
	count = 0
	batchFunc := func(items []int) error {
		count += 1
		return nil
	}

	cfg := sb.BatchConfig[int]{
		Cap:           10,
		FlushInterval: &i,
		OnFlush:       batchFunc,
	}
	b, err := sb.NewBatch(cfg)
	if err != nil {
		t.Fatalf("Failed to setup batch: %s", err)
		return nil
	}
	return b
}

// Teardown the test
func teardown(b *sb.Batch[int], t *testing.T) {
	err := b.Shutdown()
	if err != nil {
		t.Fatalf("Failed to teardown batch: %s", err)
	}
}
