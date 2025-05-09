package queue

import (
	"sync"
	"time"
)

type Task interface {
	Exec()
}

type Queue struct {
	tasks   map[time.Time][]Task
	mu      sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		tasks: make(map[time.Time][]Task),
	}
}

func (q *Queue) Add(task Task, t time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks[t] = append(q.tasks[t], task)
} 
