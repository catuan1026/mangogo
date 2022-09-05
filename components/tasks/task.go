package tasks

// TaskInf 循环执行的任务
type TaskInf interface {
	Next() bool //是否继续执行
	Run() error
}

// JobInf 时间点触发的任务 只执行一次
type JobInf interface {
	EndpointTime() int64 //触发时间
	Run() error
}
