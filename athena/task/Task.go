package task

import "sync"

type TFunc func(...any)

// Executor TaskExecutor 协程任务
type Executor struct {
	fn       TFunc
	params   []any
	callback func()
}

func NewTaskExecutor(fn TFunc, params []any, callback func()) *Executor {
	return &Executor{fn: fn, params: params, callback: callback}
}

var taskList chan *Executor // 任务列表
var taskListOnce sync.Once

func getTaskList() chan *Executor {
	taskListOnce.Do(func() {
		taskList = make(chan *Executor)
	})
	return taskList
}

// 执行任务
func (this *Executor) exec() {
	this.fn(this.params...)
}

func doTask(t *Executor) {
	go func() {
		defer func() {
			if t.callback != nil {
				t.callback()
			}
		}()
		t.exec()
	}()
}

func init() {
	getChan := getTaskList()
	go func() {
		for t := range getChan {
			doTask(t)
		}
	}()
}

// Async 异步协程任务投递
func Async(fn TFunc, callback func(), params ...any) {
	if fn == nil {
		return
	}

	go func() {
		getTaskList() <- NewTaskExecutor(fn, params, callback)
	}()
}
