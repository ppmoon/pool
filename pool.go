package pool

import "sync"

//创建worker，每一个worker抽象成一个可以执行任务的函数
type Worker struct {
	f func() error
}
//通过NewTask来创建一个worker
func NewTask(f func() error)  *Worker {
	return &Worker{
		f:f,
	}
}
//执行worker
func (t *Worker) Run(wg *sync.WaitGroup)  {
	t.f()
	//减少waitGroup计数器的值
	wg.Done()
}
//池
type Pool struct {
	//这个*Worker指针切片用来接受任务，方便外部调用，减少channel异常的问题，这里会整个切片一起提交
	//Workers []*Worker
	//这里的Worker是一个管道，用来接受其他go程带来的数据，实时执行，无限等待数据循环，这里使用另外一个管道还可以隐藏wg的操作。让外部程序使用更方便一些。
	Worker chan *Worker
	//size用来表明池的大小，不能超发。
	size int
	//jobs表示执行任务的通道用于作为队列，我们将任务从切片当中取出来，然后存放到通道当中，再从通道当中取出任务并执行。
	Jobs chan *Worker
	//用于阻塞
	wg sync.WaitGroup
}
//实例化工作池使用
func NewPool(cap int) *Pool {
	return &Pool{
		//Workers:tasks,
		Worker:make(chan *Worker),
		size: cap,
		Jobs:make(chan *Worker),
	}
}
//从jobs当中取出任务并执行。
func (p *Pool) work()  {
	for task := range p.Jobs{
		task.Run(&p.wg)
	}
}
//执行工作池当中的任务
func (p *Pool) Run(){
	//只启动有限大小的协程，协程的数量不可以超过工作池设定的数量，防止计算资源崩溃
	for i:=0;i<p.size;i++{
		go p.work()
	}
	//从worker切片当中把任务取出
	for task := range p.Worker{
		p.wg.Add(1)
		p.Jobs <- task
	}
	//执行完毕就需要关闭jobs
	close(p.Jobs)
	//执行的过程需要阻塞直到有空闲的goroutine可用
	p.wg.Wait()
}