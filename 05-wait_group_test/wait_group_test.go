package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func RunAsynchronous(group *sync.WaitGroup, in int) {
	defer group.Done()

	group.Add(1)

	fmt.Println("Run Asynchronous", in)
	time.Sleep(1 * time.Second)
}

func TestWaitGroup(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 100; i++ {
		go RunAsynchronous(&waitGroup, i)
	}

	waitGroup.Wait() // jadi ga perlu lagi pake time.Sleep untuk menunggu semua goroutine selesai, karena WaitGroup sudah menangani itu
	fmt.Println("Done")
}
