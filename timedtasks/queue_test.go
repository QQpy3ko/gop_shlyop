package queue

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	q := NewQueue()

	done := make(chan bool)
	task := &mockTask{done: done}
	timestamp := time.Now().Add(time.Second)

	q.Add(task, timestamp)

	go q.Run()

	<-done

	if !task.executed {
		t.Errorf("Expected task to be executed, but it was not")
	}
}

type mockTask struct {
	executed bool
	done     chan bool
}

func (m *mockTask) Exec() {
	m.executed = true
	m.done <- true
} 