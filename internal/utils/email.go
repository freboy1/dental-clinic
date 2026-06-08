package utils

import (
	"fmt"

	"dental_clinic/internal/config"

	"github.com/resend/resend-go/v3"
)

func SendVerificationEmail(cfx *config.Config, to, token string) error {
	client := resend.NewClient(cfx.ResendAPIKey)

	verifyLink := fmt.Sprintf(
		"%s/verify-email?token=%s",
		cfx.FrontendURL,
		token,
	)

	params := &resend.SendEmailRequest{
		From:    "Dental Clinic <onboarding@resend.dev>",
		To:      []string{to},
		Subject: "Confirm your account",
		Html: fmt.Sprintf(`
			<h2>Dental Clinic</h2>
			<p>Please confirm your email address.</p>
			<p>
				<a href="%s">
					Confirm Email
				</a>
			</p>
		`, verifyLink),
	}

	_, err := client.Emails.Send(params)
	return err
}

func SendDoctorWelcomeEmail(cfx *config.Config, to, name, confirmationCode string) error {
	client := resend.NewClient(cfx.ResendAPIKey)

	params := &resend.SendEmailRequest{
		From:    "Dental Clinic <onboarding@resend.dev>",
		To:      []string{to},
		Subject: "Welcome to Dental Clinic",
		Html: fmt.Sprintf(`
			<h2>Welcome, %s!</h2>
			<p>Your account has been created by the administrator.</p>
			<p><strong>Confirmation code:</strong> %s</p>
		`, name, confirmationCode),
	}

	_, err := client.Emails.Send(params)
	return err
}

func SendEmail(cfx *config.Config, to, subject, message string) error {
	client := resend.NewClient(cfx.ResendAPIKey)

	params := &resend.SendEmailRequest{
		From:    "Dental Clinic <onboarding@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    message,
	}

	_, err := client.Emails.Send(params)
	return err
}
