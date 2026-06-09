package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []int
	limit int
}

func NewQueue(limit int) *Queue {
	q := &Queue{limit: limit}
	q.cond = sync.NewCond(&q.mu) // 绑定互斥锁
	return q
}

// 消费者：获取数据
func (q *Queue) Pop(id int) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 必须使用 for 循环检查条件，而不是 if！
	for len(q.items) == 0 {
		fmt.Printf("消费者 %d：队列为空，开始等待...\n", id)
		q.cond.Wait() // 释放锁，进入休眠；被唤醒后自动重新加锁
	}

	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("消费者 %d：消费了数据 %d\n", id, item)

	// 消费了一个数据，队列不满完毕，通知可能在等待的生产者
	q.cond.Signal()
	return item
}

// 生产者：添加数据
func (q *Queue) Push(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 队列满了，等待
	for len(q.items) == q.limit {
		fmt.Println("生产者：队列满了，开始等待...")
		q.cond.Wait()
	}

	q.items = append(q.items, item)
	fmt.Printf("生产者：生产了数据 %d\n", item)

	// 生产了数据，通知可能在等待的消费者（这里用 Broadcast 唤醒所有消费者）
	q.cond.Broadcast()
}

func TestCond(t *testing.T) {
	q := NewQueue(2) // 容量为 2 的队列

	// 启动两个消费者
	go q.Pop(1)
	go q.Pop(2)

	time.Sleep(time.Millisecond * 100)

	// 生产者开始生产
	q.Push(10)
	time.Sleep(time.Millisecond * 100)
	q.Push(20)
	time.Sleep(time.Millisecond * 100)
	q.Push(30)

	time.Sleep(time.Second * 10)
}
