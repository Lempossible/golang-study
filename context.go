package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func parentCtx(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("parentCtx start ")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println("parentCtx done")
	defer cancel()
}

func childCtx(parentCtx context.Context) {
	for {
		select {
		case <-parentCtx.Done():
			fmt.Println("child task, parentCtx done")
			return
		case <-time.After(1000 * time.Millisecond):
			fmt.Println("child task, 1s timeout")
		}
	}
}

func main() {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	childCtx1, _ := context.WithCancel(rootCtx)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		parentCtx(rootCtx, rootCancel)
	}()
	go func() {
		defer wg.Done()
		childCtx(childCtx1)
	}()
	wg.Wait()
}
