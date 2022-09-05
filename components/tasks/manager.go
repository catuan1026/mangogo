package tasks

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskManager struct {
	taskChan chan TaskInf
}

func NewTaskManage() *TaskManager {
	return &TaskManager{
		taskChan: make(chan TaskInf, 10),
	}
}

func (t *TaskManager) AddTask(task ...TaskInf) {
	for _, v := range task {
		t.taskChan <- v
	}
}

func (t *TaskManager) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-t.taskChan: //会阻塞
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logrus.Error("task run panic:", err)
					}
				}()
				if task.EndpointTime() <= time.Now().Unix() && task.Next() {
					//执行任务
					if err := task.Run(); err != nil {
						logrus.Error("task run error:", err)
					}
				}
				if task.Next() {
					time.Sleep(time.Second * 1)
					t.taskChan <- task
				}
			}()
		}
	}
}

type JobManager struct {
	jobChan chan JobInf
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobChan: make(chan JobInf, 10),
	}
}

func (j *JobManager) AddJob(job ...JobInf) {
	for _, v := range job {
		j.jobChan <- v
	}
}

func (j *JobManager) Run(ctx context.Context) {
	/***
	ticker:=time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C: //存在溢出
			//每隔一秒检查一次
			for {
				select {
				case job := <-j.jobChan: //会阻塞
					go func() {
						defer func() {
							if err := recover(); err != nil {
								logrus.Error("job run panic:", err)
							}
						}()
						if job.EndpointTime() <= time.Now().Unix() {
							//执行任务
							if err := job.Run(); err != nil {
								logrus.Error("job run error:", err)
							}
						}else{
							//还未到执行时间
							j.jobChan <- job
						}
					}()
				default:
					return
				}
			}
		}
	}**/

	for {
		select {
		case <-ctx.Done():
			return
		case job := <-j.jobChan: //会阻塞
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logrus.Error("job run panic:", err)
					}
				}()
				if job.EndpointTime() <= time.Now().Unix() {
					//执行任务
					if err := job.Run(); err != nil {
						logrus.Error("job run error:", err)
					}
				} else {
					//还未到执行时间
					time.Sleep(time.Second * 1)
					j.jobChan <- job
				}
			}()
		}
	}
}
