package utils

import (
	"fmt"
	"net/smtp"

	"dental_clinic/internal/config"
)

func SendVerificationEmail(cfx *config.Config, to, token string) error {

	from := cfx.SMTPUser
	password := cfx.SMTPPass
	smtpHost := cfx.SMTPHost
	smtpPort := cfx.SMTPPort

	verifyLink := fmt.Sprintf(
		"http://161.35.116.104:8080/api/verify?token=%s",
		token,
	)
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

func SendDoctorWelcomeEmail(cfx *config.Config, to, name, confirmationCode string) error {
	from := cfx.SMTPUser
	password := cfx.SMTPPass
	smtpHost := cfx.SMTPHost
	smtpPort := cfx.SMTPPort

	subject := "Welcome to Dental Clinic — Your Confirmation Code\n"
	body := fmt.Sprintf(
		"Hello, %s!\n\nYour account has been created by the administrator.\n\nYour confirmation code: %s\n\nPlease use this code to confirm your account on first login.\n",
		name, confirmationCode,
	)
	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(cfx *config.Config, to, subject, message string) error {
	from := cfx.SMTPUser
	password := cfx.SMTPPass
	smtpHost := cfx.SMTPHost
	smtpPort := cfx.SMTPPort

	msg := []byte(subject + "\n" + message)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
