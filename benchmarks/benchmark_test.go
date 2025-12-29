package benchmarks

import (
	"os"
	"testing"

	sb "github.com/PeterOlsen1/superbatch"
)

var taskCount int = 100
var file *os.File

func openFile() {
	f, err := os.Open("./temp.txt")
	if err != nil {
		return
	}
	file = f
}

func closeFile() {
	if file != nil {
		file.Close()
		os.Remove("./temp.text")
	}
}

func dummyTask(i int) error {
	_, err := file.WriteString("hello, world")
	return err
}

// BenchmarkBatch-16    	  263919	      5040 ns/op	    8199 B/op	       0 allocs/op
func BenchmarkBatch(b *testing.B) {
	openFile()
	defer closeFile()

	cfg := sb.BatchConfig[int]{
		Cap:     100,
		OnFlush: dummyTask,
	}
	batch, _ := sb.NewBatch(cfg)
	defer batch.Shutdown()

	for b.Loop() {
		for range taskCount {
			batch.Add(1)
		}
	}
}

// BenchmarkGoroutines-16    	   30326	     41337 ns/op	    1719 B/op	     101 allocs/op
func BenchmarkGoroutines(b *testing.B) {
	openFile()
	defer closeFile()

	for b.Loop() {
		done := make(chan struct{}, taskCount)
		for range taskCount {
			go func() {
				dummyTask(1)
				done <- struct{}{}
			}()
		}
		for range taskCount {
			<-done
		}
	}
}

// BenchmarkSequential-16    	   77817	     14602 ns/op	   20801 B/op	     300 allocs/op
func BenchmarkSequential(b *testing.B) {
	openFile()
	defer closeFile()

	for b.Loop() {
		for range taskCount {
			dummyTask(1)
		}
	}
}
