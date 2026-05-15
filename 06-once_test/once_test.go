package main

import (
	"fmt"
	"sync"
	"testing"
)

var counter = 0

func OnlyOnce() {
	counter++
}

// TestOnce akan menjalankan fungsi OnlyOnce sebanyak 100 kali, namun karena kita menggunakan sync.Once,
// maka fungsi OnlyOnce hanya akan dieksekusi sekali saja. Sehingga nilai counter akan tetap 1
// meskipun kita memanggilnya sebanyak 100 kali.
func TestOnce(t *testing.T) {
	var once sync.Once
	var group sync.WaitGroup

	for i := 0; i < 100; i++ {
		go func() {
			group.Add(1)
			once.Do(OnlyOnce)
			group.Done()
		}()
	}

	group.Wait()
	fmt.Println("Counter:", counter)
}
