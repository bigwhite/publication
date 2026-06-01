# Go 并发竞态条件修复方案

## 问题概述
原始代码存在并发竞态条件问题，多个 goroutine 同时修改变量 `c++` 导致数据不一致。

## 检测方法
使用 Go 的 race detector 检测竞态条件：
```bash
go run -race main.go
```

## 修复方案

### 方案一：互斥锁 (sync.Mutex)
```go
var mu sync.Mutex
var c int64

for i := 0; i < int(count); i++ {
    go func() {
        mu.Lock()
        c++
        mu.Unlock()
    }()
}
```

**优点**：
- 简单直观
- 适用于复杂场景
- 代码可读性好

**缺点**：
- 性能相对较低
- 可能导致死锁

### 方案二：原子操作 (sync/atomic)
```go
var c int64

for i := 0; i < int(count); i++ {
    go func() {
        atomic.AddInt64(&c, 1)
    }()
}
```

**优点**：
- 性能较好
- 无锁设计
- 适合简单原子操作

**缺点**：
- 只适用于简单操作
- 复杂逻辑需要多个原子操作

### 方案三：原子操作 AddInt64
与方案二相同，使用 `atomic.AddInt64` 函数。

### 方案四：CAS (CompareAndSwap) 操作
```go
var c int64

for i := 0; i < int(count); i++ {
    go func() {
        for {
            current := atomic.LoadInt64(&c)
            if atomic.CompareAndSwapInt64(&c, current, current+1) {
                break
            }
        }
    }()
}
```

**优点**：
- 无锁设计
- 适合复杂原子操作
- 高性能

**缺点**：
- 代码复杂
- 可能存在自旋等待

## 性能对比
基于 100 万次操作的性能测试结果：
- 互斥锁方案: 273.8ms
- 原子操作方案: 254.5ms  
- CAS 方案: 254.9ms

## 测试结果
- **原始竞态条件代码**：结果不准确（942/1000），race detector 检测到数据竞争
- **互斥锁方案**：结果准确（1000/1000），无数据竞争
- **原子操作方案**：结果准确（1000/1000），无数据竞争
- **原子操作 AddInt64 方案**：结果准确（1000/1000），无数据竞争
- **CAS 方案**：结果准确（1000/1000），无数据竞争

## 推荐方案
1. **简单场景**：推荐使用 `atomic.AddInt64`，性能最好，代码简洁
2. **复杂场景**：推荐使用 `sync.Mutex`，代码可读性好
3. **高性能需求**：推荐使用 CAS 操作，但要注意代码复杂度

## 使用方法
1. 运行竞态条件检测：
   ```bash
   go run -race demo.go
   ```

2. 运行性能测试：
   ```bash
   go run benchmark.go
   ```

3. 运行并发安全测试：
   ```bash
   go run demo.go
   ```