package queue

import (
	"testing"
	"time"
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

type mockTask struct{}

func (m *mockTask) Exec() {

} 