package pool

import "sync"

type Worker struct {
	f func() error
}

func NewTask(f func() error)  *Worker {
	return &Worker{
		f:f,
	}
}

func (t *Worker) Run(wg *sync.WaitGroup)  {
	t.f()
	wg.Done()
}

type Pool struct {
	Workers []*Worker
	size int
	Jobs chan *Worker
	wg sync.WaitGroup
}

func NewPool(tasks []*Worker,cap int) *Pool {
	return &Pool{
		Workers:tasks,
		size: cap,
		Jobs:make(chan *Worker),
	}
}
func (p *Pool) work()  {
	for task := range p.Jobs{
		task.Run(&p.wg)
	}
}
func (p *Pool) Run(){
	for i:=0;i<p.size;i++{
		go p.work()
	}
	for _,task := range p.Workers{
		p.wg.Add(1)
		p.Jobs <- task
	}
	close(p.Jobs)
	p.wg.Wait()
}