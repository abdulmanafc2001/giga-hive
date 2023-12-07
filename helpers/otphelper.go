package helpers

import (
	"math/rand"
	"net/smtp"
	"os"
)

func GenerateOtp() int {
	return rand.Intn(9000) + 1000
}

func SendOtp(otp, email string) error {
	auth := smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("PASSWORD"), "smtp.gmail.com")
	to := []string{email}
	message := "Subject: Otp verification\nyour verification otp is " + otp
	return smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("EMAIL"), to, []byte(message))
}
