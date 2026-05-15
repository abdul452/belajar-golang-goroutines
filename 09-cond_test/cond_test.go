package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var locker = sync.Mutex{}
var cond = sync.NewCond(&locker)
var wg = sync.WaitGroup{}

func WaitCondition(value int) {
	defer wg.Done()
	// wg.Add(1) // kata gemini // ❌ Kurang tepat ditaruh di sini

	cond.L.Lock()
	cond.Wait()
	fmt.Println("Done", value)
	cond.L.Unlock()
}

func TestCond(t *testing.T) {
	// Jalankan 10 Goroutine untuk mengantre/tidur
	for i := 0; i < 10; i++ {
		wg.Add(1) // 🟢 Lebih tepat ditaruh di sini, karena kita tahu pasti akan membuat 10 Goroutine
		go WaitCondition(i)
	}

	fmt.Println("Waiting for all goroutines to be ready")
	// Beri waktu sejenak agar 10 Goroutine di atas PASTI sudah masuk posisi cond.Wait()
	time.Sleep(1 * time.Second)

	fmt.Println("All goroutines are ready, broadcasting condition")

	// menggunakan signal
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		time.Sleep(1 * time.Second)
	// 		cond.Signal()
	// 	}
	// }()

	// menggunakan broadcast
	go func() {
		cond.Broadcast()
	}()

	// 4. Tunggu sampai semua Goroutine selesai memproses sisa kodenya
	wg.Wait()
}

// dari gemini
// Fungsi pekerja biasa yang menerima pointer ke sync.WaitGroup
func pekerja(id int, wg *sync.WaitGroup) {
	// 2. Wajib panggil Done() di akhir fungsi menggunakan defer
	defer wg.Done()

	fmt.Printf("Pekerja %d: Mulai bekerja...\n", id)
	time.Sleep(1 * time.Second) // Simulasi proses kerja
	fmt.Printf("Pekerja %d: Selesai!\n", id)
}

// Fungsi Unit Test untuk menguji WaitGroup
func TestPekerja(t *testing.T) {
	var wg sync.WaitGroup
	banyakPekerja := 10

	fmt.Printf("Memulai pengujian dengan %d pekerja\n", banyakPekerja)

	for i := 1; i <= banyakPekerja; i++ {
		// 1. TAMBAH COUNTER DI SINI (Sebelum kata kunci 'go')
		wg.Add(1)

		go pekerja(i, &wg)
	}

	fmt.Println("[Main Test]: Menunggu semua pekerja selesai di latar belakang...")

	// 3. TAHAN FUNGSI TEST DI SINI sampai counter kembali ke 0
	wg.Wait()

	fmt.Println("[Main Test]: Semua pekerja beres, pengujian selesai dengan sukses!")
}

// dari gemini
// JalankanKondisi adalah fungsi biasa yang berisi logika sync.Cond
func JalankanKondisi(banyakGoroutine int) {
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)
	var wg sync.WaitGroup

	kondisiTerpenuhi := false

	// Membuat Goroutine yang menunggu sesuai jumlah parameter
	for i := 1; i <= banyakGoroutine; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			mutex.Lock()
			fmt.Println("cek kondisi")
			// ❌ SETELAH BANGUN, PROGRAM TIDAK KEMBALI KE SINI
			// Jadi baris "cek kondisi" di atas loop TIDAK AKAN PERNAH dipanggil lagi.
			for !kondisiTerpenuhi {
				fmt.Printf("Goroutine %d: Kondisi belum siap, saya tidur dulu...\n", id)
				cond.Wait() // 🛌 Goroutine tidur di sini.
				// setelah bangun, langsung cek kondisi lagi, kalau belum terpenuhi ya tidur lagi, kalau sudah terpenuhi ya lanjut eksekusi kode di bawah
				fmt.Printf("Goroutine %d: Bangun! Tapi saya cek dulu kondisi lagi...\n", id)
			}

			// jadi kenapa tidak pakai if? KLIHATANNYA BENAR, TAPI BERBAHAYA ❌
			// if !kondisiTerpenuhi {
			// 	cond.Wait()
			// }
			// Begitu bangun langsung eksekusi kode di bawah tanpa cek ulang

			fmt.Printf("Goroutine %d: Mantap! Saya sudah bangun dan koding...\n", id)
			mutex.Unlock()
		}(i)
	}

	// Beri jeda 2 detik sebelum mengubah kondisi
	time.Sleep(2 * time.Second)

	mutex.Lock()
	fmt.Println("\n[Main Program]: Mengubah kondisi menjadi TRUE...")
	kondisiTerpenuhi = false
	mutex.Unlock()

	// Bangunkan semua Goroutine
	cond.Broadcast()

	wg.Wait()
	fmt.Println("[Main Program]: Semua Goroutine selesai.")
}

func TestJalankanKondisi(t *testing.T) {
	// Kita test jalankan fungsi dengan memicu 3 Goroutine
	banyakGoroutine := 3

	t.Logf("Memulai pengujian sync.Cond dengan %d Goroutine", banyakGoroutine)

	// Memanggil fungsi biasa yang sudah kita buat
	JalankanKondisi(banyakGoroutine)

	t.Log("Pengujian sync.Cond berhasil diselesaikan tanpa deadlock")
}
