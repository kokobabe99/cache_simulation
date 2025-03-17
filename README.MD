# Cache Simulator

A cache simulator implemented in Go with Fyne GUI framework, designed to demonstrate and analyze cache access patterns.

## Features

- 4-way set-associative cache structure
- 8 sets with 4 blocks per set
- LRU (Least Recently Used) replacement policy
- Real-time visualization of cache access process
- Multiple test sequence patterns
- Detailed performance statistics

## Test Modes

1. **Sequential Test**

   - Access memory addresses in sequence
   - Tests basic cache functionality

2. **Random Test**

   - Generates random access sequences
   - Simulates real program random access patterns

3. **Mid-Repeat Test**

   - Repeats access in the middle portion of the sequence
   - Tests cache locality principle

## Cache Parameters

- Cache Access Time: 1ns
- Main Memory Access Time: 10ns
- Memory Access Penalty: 162ns (1 + 16\*10 + 1)
- Number of Sets: 8
- Ways per Set: 4

## Performance Metrics

- Access Count
- Hit Count
- Miss Count
- Hit Rate
- Miss Rate
- Average Access Time
- Total Access Time

```

h is the hit rate
C is cache access time
M is the miss penalty


T(avg) = Hit Count * C + (1-h) * miss penalty

```

# Interface Guide

![Cache Simulator Interface](docs/gui.png)

1. **Parameter Settings**

   - Configure number of cache blocks

2. **Test Options**

   - Select different test modes
   - Start simulation

3. **Visualization Area for tracing**

   - Real-time access sequence display
   - Hit/Miss status indication
   - SET index display
   - Progress bar for test completion

4. **Results Area**
   - Detailed performance statistics
   - Basic stats and performance metrics

## Development Environment

- Go 1.21.3
- Fyne v2.x

## Cache Structure

```go
type Block struct {
    SeqNumber int  // Sequence number
    Age       int  // Access age
    Valid     bool // Valid bit
}

type Set struct {
    Blocks [4]Block // 4 blocks per set
}

type Cache struct {
    Sets        [8]Set // 8 sets
    AccessCount int    // Total access count
    HitCount    int    // Cache hit count
    MissCount   int    // Cache miss count
}
```
