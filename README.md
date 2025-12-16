# Superbatch 🦸‍♂️

This is a project that was inspired by my previous load balancing project, where I needed to create a mini-batch processing library. I decided that it was important enough to spin off into its own project, leading us here.

### Initialization

```go
import batch "github.com/PeterOlsen1/superbatch"

func main() {
    batchSize := 10
    timeout := 5 * time.Millisecond
    batchFunc := func(items []int) error {
        // your function logic here...
        return nil
    }
    
    b := batch.InitBatch(batchSize, &timeout, batchFunc)
}
```

Notably, the timeout variable is expressed as a pointer becuase it can be passed in as `nil`. If timeout is nil, the batch will _only_ flush when it reaches capacity, effectively removing the minimum time requirement.

```go
b := batch.InitBatch(batchSize, nil, batchFunc)
```

### Customization

The superbatch can be customized *after* initialization for use-cases where a dynamic capacity or timeout may be necessary. In any case where a parameter is updated to be smaller than previous (capacity decrease, timeout decrease), and the current length or timer respectively is past the new value, the batch will be automatically flushed for the new parameters to work properly. See below:

```go
b := batch.InitBatch(15, &timeout, batchFunc)

for range 10 {
    b.Add(1)
}

b.SetCap(5) // this will flush, since the current batch length is 10
```

```go
timeout := 10 * time.Millisecond
b := batch.InitBatch(batchSize, &timeout, batchFunc)

time.Sleep(8 * time.Millisecond)

newTimeout := 5 * time.Millisecond
b.SetTimeout(&newTimeout) // this will flush, since the current timeout is at 5 seconds
```

(i need to actually implement the timeout one)