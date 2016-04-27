package scheduler

type TaskDetail struct {
	obj interface{}
	fun  func(name interface{})
}

func NewTask(obj interface{}, fun func(obj interface{})) *TaskDetail {
	return &TaskDetail{obj: obj, fun: fun}
}

func (j *TaskDetail) Run() {
	j.fun(j.obj)
}
