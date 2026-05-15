package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	// pool := &sync.Pool{}
	pool := &sync.Pool{
		// untuk mengisi data ketika pool kosong, kita bisa menggunakan fungsi New
		New: func() interface{} {
			return "New Data"
		},
	}
	pool.Put("Aira")
	pool.Put("Almahyra")
	pool.Put("Falah")

	for i := 0; i < 10; i++ {
		go func() {
			data := pool.Get() // Mengambil data dari pool, setelah data diambil, data tersebut akan hilang dari pool
			fmt.Println(data)
			time.Sleep(1 * time.Second) // Simulasi penggunaan data selama 1 detik
			pool.Put(data)              // Mengembalikan data ke pool setelah digunakan
		}()
	}

	time.Sleep(11 * time.Second) // Memberi waktu untuk goroutine selesai sebelum program berakhir
	fmt.Println("Selesai")
}
