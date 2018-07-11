package pool_test

import (
	"testing"
	"pool"
	"fmt"
	"time"
)
//这里使用了一个很笨的办法期望通过时间来观察goroutine开启的情况
func TestNewPool(t *testing.T) {
	w := pool.NewTask(func() error { fmt.Println(time.Now());return nil })
	//tasks := []*pool.Worker{}
	//for i:=0;i<10 ;i++  {
	//	tasks = append(tasks,w)
	//}
	//p := pool.NewPool(tasks,3)

	p := pool.NewPool(3)
	//这里启用另外一个goroutine向worker当中写入，不然会出现all goroutines are asleep，需要从管道中获得一个数据，而这个数据必须是其他goroutine线放入管道的
	go func() {
		for {
			p.Worker <- w
		}
	}()
	p.Run()
}

func TestSe(t *testing.T) {
	a := make(chan int)
	var i int
	go func() {
		for {
			i++
			a <- i
		}
	}()
	for {
		fmt.Println(<-a)
	}
}