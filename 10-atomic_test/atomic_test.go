package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// Gunakan Atomic jika yang diutak-atik hanya variabel angka tunggal,
// dan gunakan Mutex jika yang mau dilindungi adalah sebuah struct atau kumpulan baris kode logika.

func TestAtomic(t *testing.T) {
	var counter int64 = 0
	var wg sync.WaitGroup

	// Jalankan 1000 Goroutine untuk menambah counter
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // Menambah counter secara atomik
		}()
	}

	wg.Wait()                        // Tunggu sampai semua Goroutine selesai
	fmt.Println("Counter:", counter) // Harusnya selalu 1000
}
