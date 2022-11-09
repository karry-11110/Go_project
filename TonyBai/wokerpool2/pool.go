package workerpool

import (
	"errors"
	"fmt"
	"sync"
)

type Task func()

var (
	ErrWorkerPoolFreed    = errors.New("workerpool freed")
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool") // workerpool中任务已满，没有空闲goroutine用于处理新任务
)

const (
	defaultCapacity = 100
	maxCapacity     = 10000
)

type pool struct {
	capacity int

	wg sync.WaitGroup

	active chan struct{}
	tasks  chan Task
	quit   chan struct{}

	preAlloc bool
	block    bool
}

func New(capacity int, opts ...Option) *pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}
	p := &pool{
		capacity: capacity,
		active:   make(chan struct{}, capacity),
		tasks:    make(chan Task, 10),
		quit:     make(chan struct{}),
		block:    true,
	}

	for _, opt := range opts {
		opt(p)
	}
	fmt.Printf("workerpool start(preAlloc =%t)\n", p.preAlloc)
	if p.preAlloc {
		// create all goroutines and send into works channel
		for i := 0; i < p.capacity; i++ {
			p.newWorker(i + 1)
			p.active <- struct{}{}
		}
	}

	go p.run()

	return p
}
func (p *pool) returnTask(t Task) {
	go func() {
		p.tasks <- t
	}()
}
func (p *pool) run() {
	index := len(p.active)
	if !p.preAlloc {
	loop: //当Worker pool满员时，这里会走到default分支，跳出该loop，进入到下面那个for循环，后面的for循环才是真正维护这个pool的主循环
		for t := range p.tasks {
			p.returnTask(t) //这里的task仅是触发了worker创建，这里是调度循环，不处理task，所以要把task扔回tasks channel，等worker启动后再处理
			select {
			case <-p.quit:
				return
			case p.active <- struct{}{}:
				index++
				p.newWorker(index)
			default:
				break loop
			}
		}
	}
	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
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
				fmt.Printf("worker[%03d]:recover panic,", index)
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
				fmt.Printf("worker[%03d]:receive a task\n", index)
				t()
			}
		}
	}()
}
func (p *pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoIdleWorkerInPool
	}

}
func (p *pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool freed(preAlloc=%t)\n", p.preAlloc)
}
