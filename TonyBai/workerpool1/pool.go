package workpool1

import (
	"errors"
	"fmt"
	"sync"
)

type Task func()

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

var (
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool") // workerpool中任务已满，没有空闲goroutine用于处理新任务
	ErrWorkerPoolFreed    = errors.New("workerpool freed")       // workerpool已终止运行
)

type pool struct {
	capacity int           //workerpool大小
	active   chan struct{} //active channel管理是否增加worker的标记，未满就增加
	tasks    chan Task     // 分配任务的channel
	wg       sync.WaitGroup
	quit     chan struct{} //通知各个worker推出的信号channel，关闭一个无缓冲通道，能让堵塞在这个通道的子协程通畅。
}

func New(capacity int) *pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}
	p := &pool{
		capacity: capacity,
		tasks:    make(chan Task),               // 数据传递
		quit:     make(chan struct{}),           //信号传递
		active:   make(chan struct{}, capacity), //数据交流
	}
	fmt.Printf("workerpool start\n")

	go p.run()
	return p
}
func (p *pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	}
}
func (p *pool) Free() {
	close(p.quit) // make sure all worker and p.run exit and schedule return error
	p.wg.Wait()
	fmt.Printf("workerpool freed\n")
}
func (p *pool) run() {
	index := 0
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}: //如果activa满了呢？,如果不需要那么多goroutine呢？
			index++
			p.newWorker(index)
		}
	}
}
func (p *pool) newWorker(index int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", index, err)
				<-p.active
			}
			p.wg.Done()
		}()
		fmt.Printf("worker[%03d]:start\n", index)
		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]:exit\n", index)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", index)
				t()
			}
		}

	}()
}
