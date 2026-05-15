package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/* note
* Pakai time.NewTimer kalau butuh kontrol penuh (bisa di-cancel/stop di tengah jalan).
* Pakai time.After kalau cuma butuh channel delay untuk dipasang di dalam blok select (paling sering dipakai untuk timeout API).
* Pakai time.AfterFunc kalau ingin mengeksekusi suatu perintah atau fungsi secara otomatis di background setelah waktu delay tercapai.
 */

func TestTimer(t *testing.T) {
	// Membuat objek timer berdurasi 5 detik
	timer := time.NewTimer(5 * time.Second)
	fmt.Println(time.Now()) // Mencetak waktu awal mula timer dibuat

	// Kunci suksesnya di sini: Program akan BERHENTI (blocking)
	// menunggu channel 'timer.C' menerima data setelah 5 detik.
	timeValue := <-timer.C
	fmt.Println(timeValue) // Mencetak waktu saat timer tersebut meledak/selesai
}

func TestAfter(t *testing.T) {
	// Langsung mengembalikan data berupa channel saja tanpa objek struct
	channel := time.After(5 * time.Second)
	fmt.Println(time.Now())

	// Menahan program sampai channel mengirimkan data setelah 5 detik
	timeValue := <-channel // ⚠️ Catatan: sama seperti di atas, hindari nama variabel 'time'
	fmt.Println(timeValue)
}

func TestAfterFunc(t *testing.T) {
	group := sync.WaitGroup{}
	group.Add(1)

	// 1. Daftarkan fungsi. Setelah 5 detik, fungsi ini otomatis dieksekusi di background
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("Timer Selesai pada:", time.Now())
		group.Done() // Memberitahu WaitGroup kalau fungsi background sudah beres
	})

	// 2. Ini akan langsung tercetak seketika tanpa menunggu 5 detik
	fmt.Println("Main Program start:", time.Now())

	// 3. Menahan test utama agar tidak bubar duluan sebelum AfterFunc dieksekusi
	group.Wait()
}

// Contoh dasar Stop
func TestContohStop(t *testing.T) {
	timer := time.NewTimer(5 * time.Second)

	go func() {
		time.Sleep(2 * time.Second)
		// User melakukan aksi, kita stop timer-nya sebelum detik ke-5
		stopBerhasil := timer.Stop()
		if stopBerhasil {
			fmt.Println("Timer berhasil dihentikan sebelum meledak!")
		}
	}()

	// Baris ini akan membuat aplikasi hang/deadlock jika tidak di-stop dengan benar,
	// atau jika tidak dibungkus dengan select case.
}

// Simulasi database status transaksi
var statusTransaksi = "PENDING"
var mu sync.Mutex

func bayarTransaksi() {
	mu.Lock()
	statusTransaksi = "PAID"
	mu.Unlock()
}

func TestRealCaseTimerTransaction(t *testing.T) {
	// 1. User membuat transaksi. Sistem memasang timer EXPIRED selama 3 detik
	waktuExpired := 3 * time.Second
	timerExpired := time.NewTimer(waktuExpired)

	fmt.Println("[Sistem]: Transaksi dibuat (PENDING). Menunggu pembayaran maks 3 detik...")

	// Kita gunakan WaitGroup agar test utama menunggu proses background selesai
	var wg sync.WaitGroup
	wg.Add(1)

	// 2. Goroutine Monitor: Bertugas memantau apakah waktu habis duluan atau sukses duluan
	go func() {
		defer wg.Done()

		// select case digunakan untuk balapan: mana channel yang merespon duluan
		select {
		case <-timerExpired.C:
			// Jika channel timer.C menerima data duluan, artinya user TELAT bayar
			mu.Lock()
			if statusTransaksi == "PENDING" {
				statusTransaksi = "CANCELLED"
				fmt.Println("[Sistem LOG]: ⏰ Waktu habis! Transaksi otomatis diubah menjadi CANCELLED.")
			}
			mu.Unlock()

		case <-time.After(5 * time.Second):
			// Safety net agar goroutine tidak selamanya menggantung jika ada hal lain
			fmt.Println("[Sistem LOG]: Monitor selesai.")
		}
	}()

	// ====================================================================
	// SIMULASI AKSI USER (Silakan coba ganti-ganti durasi sleep di bawah ini)
	// ====================================================================

	// Skenario 1: User gerak cepat, berhasil bayar di detik ke-1 (Sebelum 3 detik)
	time.Sleep(1 * time.Second)

	fmt.Println("[User]: Melakukan pembayaran...")
	bayarTransaksi() // Mengubah status menjadi PAID

	// KUNCI UTAMA REAL CASE:
	// Karena sudah dibayar, kita WAJIB stop timer expired-nya agar transaksi tidak berubah jadi CANCELLED
	if timerExpired.Stop() {
		fmt.Println("[Sistem]: Timer expired BERHASIL dimatikan karena user sudah bayar.")
	} else {
		fmt.Println("[Sistem]: Gagal mematikan timer (mungkin waktu sudah keburu habis).")
	}

	// ====================================================================
	// SIMULASI RESET (Fitur Tambahan Waktu / Perpanjang Masa Aktif)
	// ====================================================================
	// Misal user minta perpanjangan waktu bayar karena ada gangguan, kita bisa reset timernya:
	// timerExpired.Reset(5 * time.Second) // Waktu dihitung ulang dari 0 lagi selama 5 detik ke depan

	wg.Wait() // Tunggu goroutine monitor menyelesaikan lognya

	// Verifikasi Hasil Akhir
	fmt.Printf("[Hasil Akhir]: Status transaksi kamu saat ini adalah: %s\n", statusTransaksi)

	if statusTransaksi == "PAID" {
		fmt.Println("TEST STATUS: SUCCESS (Sistem berjalan normal)")
	}
}
