package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// Task 任务结构
type Task struct {
	ID        string
	ExecuteAt time.Time // 执行时间
	Payload   string
	index     int // 在堆中的索引，用于 Update/Remove 操作（本例暂不演示 Update）
}

// TaskQueue 任务优先级队列（小顶堆）
type TaskQueue []*Task

func (pq TaskQueue) Len() int { return len(pq) }

// Less 核心逻辑：谁的时间早，谁排前面
func (pq TaskQueue) Less(i, j int) bool {
	return pq[i].ExecuteAt.Before(pq[j].ExecuteAt)
}

func (pq TaskQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *TaskQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Task)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *TaskQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // 避免内存泄漏
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Scheduler 调度器
type Scheduler struct {
	pq     TaskQueue
	mu     sync.Mutex
	wakeup chan struct{} // 唤醒信号
	stop   chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		pq:     make(TaskQueue, 0),
		wakeup: make(chan struct{}, 1),
		stop:   make(chan struct{}),
	}
}

// AddTask 添加任务
func (s *Scheduler) AddTask(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	heap.Push(&s.pq, t)

	// 如果新添加的任务是堆顶（最早要执行的），说明调度器原本等待的时间太长了
	// 需要唤醒它，让它重新计算睡眠时间
	if s.pq[0] == t {
		// 非阻塞发送
		select {
		case s.wakeup <- struct{}{}:
		default:
		}
	}
}

// Start 启动调度循环
func (s *Scheduler) Start() {
	go func() {
		for {
			var timer *time.Timer
			var d time.Duration

			s.mu.Lock()
			if s.pq.Len() == 0 {
				// 没有任务，无限等待（这里用一个极长时间代替）
				d = time.Hour * 1000
			} else {
				now := time.Now()
				top := s.pq[0]
				if top.ExecuteAt.After(now) {
					// 还没到期，计算等待时间
					d = top.ExecuteAt.Sub(now)
				} else {
					// 已经到期（或过期），立即执行
					d = 0
				}
			}
			s.mu.Unlock()

			// 创建定时器
			timer = time.NewTimer(d)

			select {
			case <-timer.C:
				// 定时器到期，说明堆顶任务该执行了
				s.mu.Lock()
				if s.pq.Len() > 0 {
					now := time.Now()
					// double check，防止并发变动
					if !s.pq[0].ExecuteAt.After(now) {
						task := heap.Pop(&s.pq).(*Task)
						s.mu.Unlock()
						// 执行任务（异步执行，避免阻塞调度器）
						go s.execute(task)
					} else {
						s.mu.Unlock()
					}
				} else {
					s.mu.Unlock()
				}

			case <-s.wakeup:
				// 被唤醒，说明有更早的任务插入了
				// 停止当前定时器，进入下一轮循环，重新计算 d
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}

			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) execute(t *Task) {
	fmt.Printf("[%s] Executing Task: %s (Scheduled: %s)\n",
		time.Now().Format("15:04:05.000"), t.Payload, t.ExecuteAt.Format("15:04:05.000"))
}

func main() {
	s := NewScheduler()
	s.Start()

	now := time.Now()
	fmt.Printf("Start Time: %s\n", now.Format("15:04:05.000"))

	// 添加任务
	s.AddTask(&Task{ID: "1", Payload: "Task 1 (Delay 2s)", ExecuteAt: now.Add(2 * time.Second)})
	s.AddTask(&Task{ID: "2", Payload: "Task 2 (Delay 5s)", ExecuteAt: now.Add(5 * time.Second)})

	// 模拟：突然插入一个更紧急的任务
	time.Sleep(1 * time.Second)
	fmt.Println("Inserting Urgent Task...")
	s.AddTask(&Task{ID: "3", Payload: "Task 3 (Delay 1.5s - Urgent)", ExecuteAt: now.Add(1500 * time.Millisecond)})

	// 阻塞主线程等待
	time.Sleep(6 * time.Second)
}
