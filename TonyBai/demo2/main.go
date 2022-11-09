package main

import (
	"fmt"
	"time"
	workerpool "wokerpool2"
)

func main() {
	p := workerpool.New(10, workerpool.WithPreAllocWorkers(false), workerpool.WithBlock(false))

	time.Sleep(2 * time.Second)
	for i := 0; i < 5; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			fmt.Printf("task[%d]: error: %s\n", i, err.Error())
		}
	}

	p.Free()
}
