package deploy

type TaskQueue struct {
	queue []TaskSpec
}

func MakeTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue: make([]TaskSpec, 10),
	}
}

func (q *TaskQueue) push(task TaskSpec) {
	q.queue = append(q.queue, task)
}

func (q *TaskQueue) pop() (TaskSpec, bool) {
	if len(q.queue) == 0 {
		return TaskSpec{}, false
	}

	task := q.queue[0]
	q.queue = q.queue[1:]
	return task, true
}

func (q *TaskQueue) isEmpty() bool {
	return len(q.queue) == 0
}
