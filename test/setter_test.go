package superbatch_test

import (
	"testing"
	"time"
)

func TestSetCap(t *testing.T) {
	b := setup(t)

	for range 9 {
		b.Add(1)
	}
	b.SetCap(8)
	if count != 9 {
		t.Errorf("Count (%d) does not match 9", count)
	}

	for range 7 {
		b.Add(1)
	}
	b.SetCap(7)
	if count != 16 {
		t.Errorf("Count (%d) does not match 16", count)
	}

	// shouldn't update, new capacity is > size
	for range 5 {
		b.Add(1)
	}
	b.SetCap(10)
	if count != 16 {
		t.Errorf("Count (%d) does not match 16", count)
	}

	teardown(b, t)
}

func TestSetInterval(t *testing.T) {
	b := setupWithInterval(t, 10*time.Millisecond)
	b.Add(1)

	time.Sleep(8 * time.Millisecond)
	newDuration := 6 * time.Millisecond
	b.SetFlushInterval(&newDuration)

	if count != 1 {
		t.Errorf("Count (%d) does not match 1", count)
	}

	time.Sleep(2 * time.Millisecond)
	newDuration = 4 * time.Millisecond
	b.SetFlushInterval(&newDuration)

	if count != 1 {
		t.Errorf("Count (%d) does not match 1", count)
	}
}
