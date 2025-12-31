package superbatch_test

import (
	"testing"
	"time"
)

// updated in batch flush functions
var count int = 0

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

func TestBatchCapacityFlush(t *testing.T) {
	b := setup(t)

	for i := range 9 {
		b.Add(1)
		if b.Len() != i+1 {
			t.Error("Batch length does not match loop iteration")
		}
	}

	b.Add(1)

	if count != 1 {
		t.Errorf("Count (%d) does not match 1", count)
	}

	teardown(b, t)
}

func TestBatchCapacityFlushTen(t *testing.T) {
	b := setup(t)

	for range 100 {
		b.Add(1)
	}

	if count != 10 {
		t.Errorf("Count (%d) does not match 10", count)
	}

	teardown(b, t)
}

func TestBatchIntervalFlush(t *testing.T) {
	b := setupWithInterval(t, 5*time.Millisecond)
	b.Add(1)
	time.Sleep(time.Nanosecond * 6000000)

	if count != 1 {
		t.Errorf("Count (%d) does not match 1", count)
	}

	teardown(b, t)
}
