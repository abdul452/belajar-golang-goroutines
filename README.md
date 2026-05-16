# Go Concurrency & Runtime Deep Dive 🚀

Repositori ini berisi kumpulan catatan eksperimen, analisis kode, dan implementasi nyata mengenai konsep *Concurrency Control* serta manajemen *Runtime* di bahasa pemrograman Go. Fokus utama dari repositori ini adalah memahami mekanisme kerja *under the hood*, menghindari jebakan umum (*anti-patterns*), dan menerapkan standar industri backend produksi.

---

## 📌 Daftar Isi Belajar

1. [Advanced Sync Cond & Spurious Wakeup](#1-advanced-sync-cond--spurious-wakeup)
2. [Best Practice sync.WaitGroup](#2-best-practice-syncwaitgroup)
3. [Performance Tuning: Atomic vs Mutex](#3-performance-tuning-atomic-vs-mutex)
4. [Asynchronous Control: Timer & Ticker Memory Management](#4-asynchronous-control-timer--ticker-memory-management)
5. [Deep Dive Go Runtime & GOMAXPROCS](#5-deep-dive-go-runtime--gomaxprocs)

---

## 1. Advanced Sync Cond & Spurious Wakeup

### Masalah & Analisis
Saat menggunakan `sync.Cond` untuk menidurkan Goroutine via `cond.Wait()`, penggunaan struktur kondisi `if` sangat berbahaya karena rentan terhadap fenomena **Spurious Wakeup** (Goroutine terbangun palsu secara tidak sengaja oleh sistem/OS meskipun sinyal belum dikirim).

### Solusi & Implementasi Loop `for`
Pola yang aman adalah membungkus pengecekan *state variable* di dalam perulangan `for`. Ketika Goroutine dibangunkan oleh `cond.Broadcast()`, alur eksekusi akan dipaksa mengecek ulang variabel kondisi sebelum diizinkan keluar dari loop.

```go
mutex.Lock()
// Setelah bangun, program TIDAK KEMBALI KE SINI
for !kondisiTerpenuhi {
    fmt.Printf("Goroutine %d: Kondisi belum siap, saya tidur dulu...\n", id)
    cond.Wait() // Goroutine tidur dan melepas lock di sini
    
    // Begitu bangun, baris ini dieksekusi, lalu melompat ke atas mengecek loop 'for' kembali
    fmt.Printf("Goroutine %d: Bangun! Tapi saya cek dulu kondisi lagi...\n", id)
}
fmt.Printf("Goroutine %d: Mantap! Saya sudah bangun dan koding...\n", id)
mutex.Unlock()
```

## 2. Best Practice sync.WaitGroup
### Aturan Emas Penempatan wg.Add()
Untuk menghindari kendala kejar-kejaran (race condition) antara program utama (main thread) dan background job, registrasi counter wg.Add(1) wajib ditempatkan di dalam fungsi utama sebelum kata kunci go dip
