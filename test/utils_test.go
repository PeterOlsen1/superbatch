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
	b, err := sb.InitBatch(10, nil, batchFunc)
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
	b, err := sb.InitBatch(10, &i, batchFunc)
	if err != nil {
		t.Fatalf("Failed to setup batch: %s", err)
		return nil
	}
	return b
}

// Teardown the test
func teardown(b *sb.Batch[int], t *testing.T) {
	err := b.Close()
	if err != nil {
		t.Fatalf("Failed to teardown batch: %s", err)
	}
}
