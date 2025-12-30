// Package superbatch provides utilities for easy and efficient batch processing in Go.
//
// This package is designed to help manage batching of data with features like:
// - Automatic flushing based on time intervals or capacity limits.
// - Thread-safe operations for concurrent use.
// - Customizable flush logic to suit various use cases.
//
// Example usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "time"
//	    "github.com/yourusername/superbatch"
//	)
//
//	func main() {
//	    onFlush := func(items []int) error {
//	        fmt.Println("Flushed items:", items)
//	        return nil
//	    }
//
//		interval := time.Second * 5
//	    batch := superbatch.NewBatch[int](10, &interval, onFlush)
//	    defer batch.Shutdown()
//
//	    batch.Add(1)
//	    batch.Add(2)
//	    // Items will be flushed automatically based on the interval or capacity.
//	}
package superbatch
