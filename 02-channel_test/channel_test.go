package main

import (
	"fmt"
	"testing"
	"time"
)

/*
* ketika kita membuat channel, kita bisa memilih untuk membuat channel dengan buffer atau tanpa buffer.
* Channel tanpa buffer hanya dapat menampung satu data pada satu waktu, sedangkan channel dengan buffer dapat menampung beberapa data sekaligus.
* Ketika kita mengirim data ke channel tanpa buffer, pengirim akan terblokir hingga penerima menerima data tersebut. Sebaliknya, ketika kita mengirim data ke channel dengan buffer, pengirim tidak akan terblokir selama buffer belum penuh. Namun, jika buffer sudah penuh, pengirim akan terblokir hingga ada ruang kosong di dalam buffer. Demikian pula, penerima akan terblokir jika mencoba menerima data dari channel yang kosong, baik itu channel dengan buffer maupun tanpa buffer.
* Selain itu, kita juga dapat menutup channel setelah selesai digunakan. Menutup channel memberi tahu penerima bahwa tidak akan ada data lagi yang dikirim melalui channel tersebut. Penerima dapat menggunakan pernyataan range untuk menerima data dari channel hingga channel ditutup. Jika kita mencoba mengirim data ke channel yang sudah ditutup, akan terjadi panic. Oleh karena itu, penting untuk memastikan bahwa kita hanya menutup channel setelah semua pengirim selesai mengirim data dan semua penerima selesai menerima data.
* Dengan memahami konsep channel, kita dapat mengelola komunikasi antar goroutine dengan lebih efektif dan menghindari masalah seperti
* deadlock atau race condition yang dapat terjadi ketika menggunakan goroutine tanpa koordinasi yang tepat. Channel memungkinkan kita untuk menyinkronkan goroutine dan memastikan bahwa data dikirim dan diterima dengan benar, sehingga meningkatkan keandalan dan kinerja aplikasi kita.
* Dalam contoh di atas, kita membuat channel dengan tipe data string dan menjalankan goroutine untuk mengirim data ke channel. Kemudian, kita menerima data dari channel dan mencetaknya. Kita juga memberikan waktu untuk goroutine selesai sebelum program utama berakhir dengan menggunakan time.Sleep. Selain itu, kita juga menunjukkan contoh penggunaan channel dengan buffer dan cara menutup channel setelah selesai digunakan. Dengan memahami konsep channel, kita dapat mengelola komunikasi antar goroutine dengan lebih efektif dan menghindari masalah seperti deadlock atau race condition yang dapat terjadi ketika menggunakan goroutine tanpa koordinasi yang tepat. Channel memungkinkan kita untuk menyinkronkan goroutine dan memastikan bahwa data dikirim dan diterima dengan benar, sehingga meningkatkan keandalan dan kinerja aplikasi kita.
 */

func TestCreateChannel(t *testing.T) {
	channel := make(chan string) // Membuat channel dengan tipe data string
	defer close(channel)         // Menutup channel setelah selesai digunakan

	go func() { // Menjalankan goroutine untuk mengirim data ke channel
		channel <- "Hello Channel"          // Mengirim data ke channel
		fmt.Println("Data sent to channel") // Output: Data sent to channel
	}()

	message := <-channel        // Menerima data dari channel
	fmt.Println(message)        // Output: Hello Channel
	time.Sleep(5 * time.Second) // Memberi waktu untuk goroutine selesai sebelum program utama berakhir
}

func TestChannelWithBuffer(t *testing.T) {
	channel := make(chan string, 2) // Membuat channel dengan buffer berkapasitas 2

	channel <- "Message 1" // Mengirim data ke channel
	channel <- "Message 2" // Mengirim data ke channel

	fmt.Println(<-channel) // Menerima data pertama dari channel
	fmt.Println(<-channel) // Menerima data kedua dari channel
}

func TestChannelClose(t *testing.T) {
	channel := make(chan string)

	go func() {
		channel <- "Hello Channel"
		close(channel) // Menutup channel setelah mengirim data
	}()

	for message := range channel { // Menerima data dari channel hingga channel ditutup
		fmt.Println(message) // Output: Hello Channel
	}
}

// channel parameter digunakan untuk mengirim data dari satu goroutine ke goroutine lainnya. Dengan menggunakan channel, kita dapat mengkoordinasikan komunikasi antar goroutine dengan lebih efektif dan menghindari masalah seperti
func GiveMeResponse(channel chan string) {
	time.Sleep(2 * time.Second)               // Simulasi proses yang memakan waktu
	channel <- "Response from GiveMeResponse" // Mengirim data ke channel
}

func TestChannelAsParameter(t *testing.T) {
	channel := make(chan string)

	go GiveMeResponse(channel) // Menjalankan goroutine dengan channel sebagai parameter

	response := <-channel // Menerima data dari channel
	fmt.Println(response) // Output: Response from GiveMeResponse
}

// channel in dan out
func OnlyIn(channel chan<- string) {
	time.Sleep(2 * time.Second)  // Simulasi proses yang memakan waktu
	channel <- "Data for OnlyIn" // Mengirim data ke channel
}

func OnlyOut(channel <-chan string) {
	response := <-channel // Menerima data dari channel
	fmt.Println(response) // Output: Data for OnlyIn
}

func TestChannelInOut(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go OnlyIn(channel)  // Menjalankan goroutine untuk mengirim data ke channel
	go OnlyOut(channel) // Menjalankan goroutine untuk menerima data dari channel

	time.Sleep(5 * time.Second) // Memberi waktu untuk goroutine selesai sebelum program utama berakhir
}

// Buffer channel
func TestBufferChannel(t *testing.T) {
	channel := make(chan string, 3) // Membuat channel dengan buffer berkapasitas 3
	defer close(channel)

	// channel <- "Message 1" // Mengirim data ke channel
	// channel <- "Message 2" // Mengirim data ke channel
	// channel <- "Message 3" // Mengirim data ke channel

	// fmt.Println(<-channel) // Menerima data pertama dari channel
	// fmt.Println(<-channel) // Menerima data kedua dari channel
	// fmt.Println(<-channel) // Menerima data ketiga dari channel

	go func() {
		channel <- "Message 1"
		channel <- "Message 2"
		channel <- "Message 3"
	}()

	go func() {
		fmt.Println(<-channel) // Menerima data pertama dari channel
		fmt.Println(<-channel) // Menerima data kedua dari channel
		fmt.Println(<-channel) // Menerima data ketiga dari channel
	}()

	fmt.Println("Selesai")
	time.Sleep(5 * time.Second) // Memberi waktu untuk goroutine selesai sebelum program utama berakhir
}

// range channel
func TestRangeChannel(t *testing.T) {
	channel := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			channel <- fmt.Sprintf("Message %d", i) // Mengirim data ke channel
		}
		close(channel) // Menutup channel setelah mengirim semua data
	}()

	for data := range channel { // Menerima data dari channel hingga channel ditutup
		fmt.Println(data) // Output: Message 0, Message 1, Message 2, ..., Message 9
	}

	fmt.Println("Selesai")
}

// select channel
func TestSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	defer close(channel1)
	defer close(channel2)

	go GiveMeResponse(channel1) // Menjalankan goroutine untuk mengirim data ke channel1
	go GiveMeResponse(channel2) // Menjalankan goroutine untuk mengirim data ke channel2

	counter := 0
	for {
		select {
		case response1 := <-channel1: // Menerima data dari channel1
			fmt.Println("Received from channel1:", response1)
			counter++
		case response2 := <-channel2: // Menerima data dari channel2
			fmt.Println("Received from channel2:", response2)
			counter++
		}
		if counter == 2 {
			break
		}
	}

	fmt.Println("Selesai")
}

func TestDeafultSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	defer close(channel1)
	defer close(channel2)

	go GiveMeResponse(channel1) // Menjalankan goroutine untuk mengirim data ke channel1
	go GiveMeResponse(channel2) // Menjalankan goroutine untuk mengirim data ke channel2

	counter := 0
	for {
		select {
		case response1 := <-channel1: // Menerima data dari channel1
			fmt.Println("Received from channel1:", response1)
			counter++
		case response2 := <-channel2: // Menerima data dari channel2
			fmt.Println("Received from channel2:", response2)
			counter++
		default:
			fmt.Println("No response received yet")
			time.Sleep(500 * time.Millisecond) // Memberi waktu sebelum mencoba lagi
		}
		if counter == 2 {
			break
		}
	}

	fmt.Println("Selesai")
}
