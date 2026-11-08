package deploy

import (
	"context"
	"runtime"
	"sync"
)

// Deploy 调度器
type Deploy struct {
	// 最大 worker 数量
	maxWorker int
	// 任务队列
	queue *TaskQueue
	// 有新任务
	queueNotify chan struct{}
	// 任务队列锁
	queueMu sync.Mutex
}

func (d *Deploy) Run() {
	for {
		<-d.queueNotify

		d.queueMu.Lock()
		task, ok := d.queue.pop()
		if ok {
			continue
		}

		ctx, _ := context.WithCancel(context.Background())
		task.FnCtx(ctx)

		d.queueMu.Unlock()
	}
}

func (d *Deploy) notice() {
	d.queueNotify <- struct{}{}
}

func (d *Deploy) AddTask(task TaskSpec) {
	d.queueMu.Lock()
	d.queue.push(task)
	d.queueMu.Unlock()
	d.notice()
}

func MakeDeploy() *Deploy {
	deploay := &Deploy{
		maxWorker: runtime.NumCPU() * 2,
		queue:     MakeTaskQueue(),
	}

	for i := 0; i < deploay.maxWorker; i++ {
		go deploay.Run()
	}

	return deploay
}
