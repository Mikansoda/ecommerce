package service

import (
	"fmt"

	"marketplace/config"
	gomail "gopkg.in/gomail.v2"
)

func SendEmailOTP(to, otp string) error {
	if config.C.Env == "dev" {
        fmt.Println("[DEV MODE] OTP untuk", to, ":", otp)
        return nil
    }
    // kalau prod, beneran kirim
	
	m := gomail.NewMessage()
	m.SetHeader("From", config.C.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain", fmt.Sprintf("Kode OTP kamu: %s (berlaku 10 menit).", otp))

	d := gomail.NewDialer(config.C.SMTPHost, config.C.SMTPPort, config.C.SMTPUser, config.C.SMTPPass)
	return d.DialAndSend(m)
}
