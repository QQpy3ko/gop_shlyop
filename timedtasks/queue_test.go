package queue

import (
	"testing"
	"time"
    "sync/atomic"
)

func TestAdd(t *testing.T) {
	q := NewQueue()

	task := &mockTask{}
	timestamp := time.Now().Add(time.Hour)

	q.Add(task, timestamp)

	if len(q.tasks[timestamp]) != 1 {
		t.Errorf("Expected 1 task in queue, but got %d", len(q.tasks[timestamp]))
	}

	if q.tasks[timestamp][0] != task {
		t.Errorf("Expected added task to be in queue, but got different task")
	}
}

func TestRun(t *testing.T) {
	q := NewQueue()

	var executed int32
	task := &mockTask{executed: &executed}
	timestamp := time.Now().Add(time.Second)

	q.Add(task, timestamp)

	go q.Run()

	time.Sleep(2 * time.Second)

	if atomic.LoadInt32(&executed) != 1 {
		t.Errorf("Expected task to be executed, but it was not")
	}
}

type mockTask struct {
	executed *int32
}

func (m *mockTask) Exec() {
	atomic.AddInt32(m.executed, 1)
} 