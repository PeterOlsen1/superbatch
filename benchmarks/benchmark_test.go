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

func dummyTask(line string) error {
	_, err := file.WriteString(line)
	return err
}

func dummyTaskBatch(lines []string) error {
	for _, s := range lines {
		_, err := file.WriteString(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// BenchmarkBatch-8   	  264058	      5449 ns/op	   17216 B/op	       0 allocs/op
func BenchmarkBatch(b *testing.B) {
	openFile()
	defer closeFile()

	cfg := sb.BatchConfig[string]{
		Cap:     100,
		OnFlush: dummyTaskBatch,
	}
	batch, _ := sb.NewBatch(cfg)
	defer batch.Shutdown()

	for b.Loop() {
		for range taskCount {
			batch.Add("hello, world!")
		}
	}
}

// BenchmarkThreadedBatch-8   	  290559	      4338 ns/op	   15646 B/op	       0 allocs/op
func BenchmarkThreadedBatch(b *testing.B) {
	openFile()
	defer closeFile()

	cfg := sb.BatchConfig[string]{
		Cap:      100,
		OnFlush:  dummyTaskBatch,
		Threaded: true,
	}
	batch, _ := sb.NewBatch(cfg)
	defer batch.Shutdown()

	for b.Loop() {
		for range taskCount {
			batch.Add("hello, world!")
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
				dummyTask("hello, world!")
				done <- struct{}{}
			}()
		}
		for range taskCount {
			<-done
		}
	}
}

// BenchmarkSequential-16    	 3826318	       314.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSequential(b *testing.B) {
	openFile()
	defer closeFile()

	for b.Loop() {
		for range taskCount {
			dummyTask("hello, world!")
		}
	}
}
