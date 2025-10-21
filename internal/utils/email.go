package utils

import (
	"dental_clinic/internal/config"
	"fmt"
	"net/smtp"
)

func SendVerificationEmail(cfx *config.Config, to, token string) error {

	from := cfx.SMTPUser
	password := cfx.SMTPPass
	smtpHost := cfx.SMTPHost
	smtpPort := cfx.SMTPPort


	verifyLink := fmt.Sprintf("http://localhost:8080/api/verify?token=%s", token)

	subject := "Confirm your account\n"
	body := fmt.Sprintf("click to activate account:\n%s", verifyLink)
	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil

}