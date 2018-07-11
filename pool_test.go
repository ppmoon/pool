package pool_test

import (
	"testing"
	"pool"
	"fmt"
	"time"
)
//这里使用了一个很笨的办法期望通过时间来观察goroutine开启的协程效果。
func TestNewPool(t *testing.T) {
	w := pool.NewTask(func() error { fmt.Println(time.Now());return nil })
	tasks := []*pool.Worker{}
	for i:=0;i<10 ;i++  {
		tasks = append(tasks,w)
	}
	p := pool.NewPool(tasks,3)
	p.Run()
}