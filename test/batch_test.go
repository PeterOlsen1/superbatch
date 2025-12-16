package superbatch_test

import (
	"testing"

	sb "github.com/PeterOlsen1/superbatch"
)

var count int = 1

// Initialize a batch of size 10 with no timeout.
// If creation fails, test fails.
func setup(t *testing.T) *sb.Batch[int] {
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

// Teardown the test
func teardown(b *sb.Batch[int], t *testing.T) {
	err := b.Close()
	if err != nil {
		t.Fatalf("Failed to teardown batch: %s", err)
	}
}

func TestCreation(t *testing.T) {
	b := setup(t)
	teardown(b, t)
}

func TestBatchAdd(t *testing.T) {
	b := setup(t)

	err := b.Add(1)
	if err != nil {
		t.Errorf("Failed to add to batch: %s", err)
	}

	teardown(b, t)
}

func TestBatchFlush(t *testing.T) {

}
