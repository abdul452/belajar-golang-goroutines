package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	// Inisialisasi Ticker dengan interval 1 detik.
	// Mulai detik ini, setiap 1 detik, Go akan mengirim data waktu ke ticker.C
	ticker := time.NewTicker(1 * time.Second)

	done := make(chan bool) // Channel bantuan untuk menyetop perulangan

	// Goroutine Penyelamat (Asinkron)
	go func() {
		// Program utama akan sibuk nge-loop di bawah.
		// Goroutine ini berjalan di background, tidur selama 5 detik...
		time.Sleep(5 * time.Second)

		// ...lalu menyetop Ticker agar tidak memompa data lagi.
		ticker.Stop() // Setop mesin tickernya
		done <- true  // Kirim sinyal ke program utama kalau durasi sudah habis
		fmt.Println("[Background]: Ticker berhasil di-stop!")
	}()

	// Perulangan Range Channel
	fmt.Println("Memulai Ticker...")

	// Menggunakan select case di dalam loop untuk menghindari deadlock
Luar: // 🎯 Ini Label-nya (harus diakhiri tanda titik dua ':')
	for {
		select {
		case <-done:
			// Jika channel done menerima data, kita keluar dari loop utama
			fmt.Println("Stop menerima sinyal ticker.")
			break Luar
			// break // ❌ SIALNYA: break ini cuma mengeluarkan kita dari blok SELECT (Tingkat 2)
		case waktu := <-ticker.C:
			// Selama belum ada sinyal done, cetak data dari ticker setiap detik
			fmt.Println("Tick pada:", waktu)
		}
	}

	// 🎯 Begitu 'break Luar' dieksekusi, program langsung lompat mendarat di sini
	fmt.Println("Test Selesai dengan sukses!")

}

func TestTick(t *testing.T) {
	channel := time.Tick(1 * time.Second)

	for waktu := range channel {
		fmt.Println(waktu)
	}
}

// dari gemini
func TestTickYangBenar(t *testing.T) {
	// 1. Gunakan NewTicker agar mesinnya bisa kita STOP nanti
	ticker := time.NewTicker(1 * time.Second)

	// Pastikan ticker selalu dimatikan saat fungsi test ini selesai (mencegah memory leak)
	defer ticker.Stop()

	// 2. Gunakan Context untuk mengatur batas waktu (Timeout) test kita selama 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("Memulai cetak data setiap detik...")

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				// Jika waktu 5 detik dari context sudah habis, keluar dari loop
				fmt.Println("[Sistem]: Waktu pengujian habis, menghentikan loop.")
				return

			case waktu := <-ticker.C:
				// Cetak waktu setiap detik dari ticker
				fmt.Println("Tick pada:", waktu.Format("15:04:05"))
			}
		}
	}()

	// Tunggu sampai Goroutine background selesai memproses loop
	wg.Wait()
	fmt.Println("Test Selesai dengan sukses dan aman dari memory leak!")
}
