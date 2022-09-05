package tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type TaskManager struct {
	Interval  time.Duration //执行间隔
	taskChain []TaskInf
	mu        sync.Mutex
}

func (t *TaskManager) AddTask(task ...TaskInf) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.taskChain = append(t.taskChain, task...)
}

func (t *TaskManager) DelTask(index int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if index < 0 || index >= len(t.taskChain) {
		return
	}
	t.taskChain = append(t.taskChain[:index], t.taskChain[index+1:]...)
}

func (t *TaskManager) Run(ctx context.Context) {
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, task := range t.taskChain {
				if task.Next() {
					go func() {
						defer func() {
							if err := recover(); err != nil {
								logrus.Error("task run panic:", err)
							}
						}()
						t := task
						if err := t.Run(); err != nil {
							logrus.Error("task run error:", err)
						}
					}()
				}
			}
		}
	}
}
