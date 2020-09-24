package taskit

import (
	"github.com/FTwOoO/util/errorkit"
	"github.com/FTwOoO/util/logging"
	"sync"
)

var (
	ErrTaskNotExitst   = errorkit.NewStructuredError().AddParam(logging.KeyEvent, "taskNotExits")
	defaultTaskManager = newTasks()
)

func Run(id string, f func() (data interface{}, err error)) {
	defaultTaskManager.Run(id, f)
}

func Wait(id string) (data interface{}, err error) {
	return defaultTaskManager.Wait(id)
}

type taskManager struct {
	m *sync.Map
}

func newTasks() *taskManager {
	return &taskManager{
		m: &sync.Map{},
	}
}

func (this *taskManager) Wait(id string) (data interface{}, err error) {
	task, ok := this.m.Load(id)
	if !ok {
		err = ErrTaskNotExitst
		return
	}

	task.(*taskGoGroutine).wait()
	data = task.(*taskGoGroutine).Data
	err = task.(*taskGoGroutine).Error
	this.m.Delete(id)
	return
}

func (this *taskManager) Run(taskId string, f func() (data interface{}, err error)) (id uint64) {
	task := newTaskGoGroutine(taskId)
	this.m.Store(taskId, task)

	go func(task *taskGoGroutine) {
		logging.Log.Infow(logging.KeyScope, "gotask", logging.KeyEvent, "taskStart", "id", task.Id)

		defer func() {
			re := recover()
			if re != nil {
				task.SetPanic(re)
			}

			task.exit()
			this.m.Delete(task.Id)
		}()

		var data interface{}
		var err error
		data, err = f()
		if err != nil {
			err = errorkit.WrapError(err).SetScope("gotask").AddParam("taskId", task.Id)
		} else {
			logging.Log.Infow(logging.KeyScope, "gotask", logging.KeyEvent, "taskEnd", "id", task.Id)
		}
		task.SetResult(data, err)

	}(task)
	return
}

type taskGoGroutine struct {
	Id string

	//这三个字段不存在并发读写的问题，因为写肯定在读之前，因为读是在
	//done之后
	Data     interface{}
	Error    error
	IsPannic bool

	done chan int
}

func newTaskGoGroutine(id string) *taskGoGroutine {
	return &taskGoGroutine{
		Id:   id,
		done: make(chan int),
	}
}

func (this *taskGoGroutine) SetPanic(e interface{}) {
	this.IsPannic = true
	this.Error = e.(error)
}

func (this *taskGoGroutine) SetResult(data interface{}, err error) {
	this.Data = data
	this.Error = err
}

func (this *taskGoGroutine) GetResult() (data interface{}, err error) {
	return this.Data, this.Error
}

func (this *taskGoGroutine) wait() {
	<-this.done
}

func (this *taskGoGroutine) exit() {
	close(this.done)
}
