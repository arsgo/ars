package scheduler

type TaskDetail struct {
	name string
	fun  func(name string)
}

func NewTask(name string, fun func(name string)) *TaskDetail {
	return &TaskDetail{name: name, fun: fun}
}

func (j *TaskDetail) Run() {
	j.fun(j.name)
}
