package gomaxprocstest

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestGetGomaxprocs(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Second)
			wg.Done()
		}()
	}
	totalCpu := runtime.NumCPU()
	fmt.Println("Total CPU", totalCpu)

	// Angka -1 di sini artinya kita cuma mau mengintip/membaca nilainya tanpa mengubahnya.
	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println("Total Thread", totalThread)

	totalGoroutine := runtime.NumGoroutine()
	fmt.Println("Total Goroutine", totalGoroutine)

	wg.Wait()
}
