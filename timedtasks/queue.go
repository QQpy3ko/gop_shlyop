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

func (q *Queue) Run() {
	for {
		now := time.Now()
		q.mu.Lock()
		for t, tasks := range q.tasks {
			if t.Before(now) || t.Equal(now) {
				for _, task := range tasks {
					go task.Exec()
				}
				delete(q.tasks, t)
			}
		}
		q.mu.Unlock()
		time.Sleep(time.Second)
	}
} 
