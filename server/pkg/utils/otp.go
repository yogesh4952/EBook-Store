package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) GenerateOTP() (string, error) {
	max := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (e *EmailService) SentOTPEmail(toEmail, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: Your Verification Code\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
body{font-family:Arial,sans-serif;background:#f4f4f4;padding:40px 0;text-align:center;}
.container{max-width:480px;margin:0 auto;background:#fff;padding:40px 30px;border-radius:12px;box-shadow:0 4px 20px rgba(0,0,0,0.08);}
h1{color:#2d3436;font-size:26px;margin-bottom:10px;}
.otp-box{font-size:36px;font-weight:bold;color:#0984e3;letter-spacing:8px;padding:20px;background:#e8f4fd;border-radius:8px;margin:25px 0;}
.subtitle{color:#636e72;font-size:15px;line-height:1.6;}
.footer{margin-top:30px;font-size:12px;color:#b2bec3;}
</style>
</head>
<body>
<div class="container">
<h1>Book Store Nepal</h1>
<p class="subtitle">Use the code below to verify your identity and log in.</p>
<div class="otp-box">%s</div>
<p class="subtitle">This code expires in 5 minutes. If you didn't request it, ignore this email.</p>
<div class="footer">Book Store Nepal &mdash; Secure Auth System</div>
</div>
</body>
</html>`, otp)

	msg := []byte(subject + mime + body)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
}
