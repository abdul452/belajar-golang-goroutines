package main

import (
	"fmt"
	"testing"
	"time"
)

/*
* Di dalam bahasa Go, ketika kita ingin memicu Goroutine, ada 3 jenis fungsi yang biasa kita gunakan berdasarkan strukturnya:
*
* 1. Fungsi Biasa (Named Function)
*
* 2. Fungsi Anonim Tanpa Parameter (Anonymous Function)
*
* 3. Fungsi Anonim Menggunakan Parameter (Anonymous Function with Parameters)
 */

func RunHelloWorld() {
	fmt.Println("Hello World")
}

// 1. Fungsi Biasa (Named Function)
func TestNamedGoroutine(t *testing.T) {
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

// 2. Fungsi Anonim Tanpa Parameter (Anonymous Function)
func TestAnonymousNoParam(t *testing.T) {

	go func() {
		// Tulis logika kodinganmu di sini
		fmt.Println("Goroutine anonim tanpa parameter jalan!")
	}() // 💡 INGAT: Tanda () di sini artinya "JALANKAN SEKARANG"

}

// 3. Fungsi Anonim Menggunakan Parameter (Anonymous Function with Parameters)
func TestAnonymousWithParam(t *testing.T) {

	for i := 1; i <= 3; i++ {
		// 1. Definisikan parameter di 'func(id int)'
		go func(id int) {
			fmt.Println("Goroutine ID:", id)
		}(i) // 2. KIRIM variabel 'i' ke dalam kurung penutup di sini!
	}

	time.Sleep(1 * time.Second) // Hanya untuk nunggu log muncul
}
