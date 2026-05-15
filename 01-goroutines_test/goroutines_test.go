package main

import (
	"fmt"
	"testing"
	"time"
)

func RunHelloWorld() {
	fmt.Println("Hello World")
}

func TestCreateGoroutine(t *testing.T) {
	go RunHelloWorld() // Menjalankan fungsi RunHelloWorld sebagai goroutine
	fmt.Println("Test Create Goroutine")

	time.Sleep(1 * time.Second) // Memberi waktu untuk goroutine selesai sebelum program utama berakhir
}

func DisplayNumber(number int) {
	fmt.Println("Number:", number)
}

func TestCreateGoroutineWithParameter(t *testing.T) {
	for i := 0; i < 100000; i++ {
		go DisplayNumber(i) // Menjalankan fungsi DisplayNumber sebagai goroutine dengan parameter
	}

	time.Sleep(10 * time.Second) // Memberi waktu untuk semua goroutine selesai sebelum program utama berakhir
}
