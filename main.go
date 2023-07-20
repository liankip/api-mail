package main

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/hibiken/asynq/v2"
)

func sendEmailTask(to, from, subject, body string) asynq.TaskResult {
	// Kredensial Mailtrap.io (ganti dengan kredensial Anda)
	host := "smtp.mailtrap.io"
	port := 587
	username := "your_mailtrap_username"
	password := "your_mailtrap_password"

	// Konfigurasi server SMTP Mailtrap.io
	auth := smtp.PlainAuth("", username, password, host)

	// Pesan email
	message := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Mengirim email melalui server Mailtrap.io
	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, []string{to}, message)
	if err != nil {
		log.Println("Failed to send email:", err)
		return asynq.NewTaskResultError("Failed to send email")
	}

	log.Println("Email sent successfully to:", to)
	return asynq.NewTaskResult(nil)
}

func main() {
	rdb := asynq.RedisClientOpt{
		Addr:     "localhost:6379", // Ganti dengan alamat Redis Anda
		Password: "",               // Kosongkan jika tidak ada password
		DB:       0,                // Gunakan DB 0 untuk asynq
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt(rdb),
		asynq.Config{
			Concurrency: 10, // Sesuaikan dengan kebutuhan Anda
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("send-email", sendEmailTask)

	go srv.Run(mux)

	client := asynq.NewClient(rdb)

	// Jadwalkan tugas pengiriman email menggunakan Asynq
	to := "recipient@example.com"
	from := "sender@example.com"
	subject := "Testing Mailtrap.io with Asynq"
	body := "This is a test email using Mailtrap.io with Asynq!"

	task := asynq.NewTask("send-email", map[string]interface{}{
		"to":      to,
		"from":    from,
		"subject": subject,
		"body":    body,
	})

	delay := 10 * time.Second
	_, err := client.Enqueue(task, asynq.ProcessIn(delay))
	if err != nil {
		log.Println("Failed to enqueue task:", err)
	}

	time.Sleep(5 * time.Second) // Tunggu agar tugas berhasil dieksekusi

	// Pastikan untuk mengganti `your_mailtrap_username` dan `your_mailtrap_password`
	// dengan kredensial Mailtrap.io Anda. Juga, sesuaikan alamat penerima dan pengirim
	// email dengan alamat yang sesuai.

	log.Println("Done!")
}
