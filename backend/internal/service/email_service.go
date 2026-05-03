package service

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	apiKey    string
	fromEmail string
}

func NewEmailService(apiKey, fromEmail string) *EmailService {
	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}
}

// SendOTPEmail sends OTP verification email
func (s *EmailService) SendOTPEmail(toEmail, otp string) error {
	if s.apiKey == "" {
		// Fallback to console logging for development
		log.Printf("OTP Email to %s: %s\n", toEmail, otp)
		return nil
	}

	from := mail.NewEmail("MasterGo", s.fromEmail)
	to := mail.NewEmail("", toEmail)

	subject := "Verify Your Email - MasterGo"
	plainTextContent := fmt.Sprintf("Your verification code is: %s. It will expire in 10 minutes.", otp)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Email Verification</h2>
				<p>Your verification code is:</p>
				<h1 style="color: #007bff;">%s</h1>
				<p>This code will expire in 10 minutes.</p>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
		</html>
	`, otp)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}

// SendPasswordResetEmail sends password reset link
func (s *EmailService) SendPasswordResetEmail(toEmail, resetLink string) error {
	if s.apiKey == "" {
		log.Printf("Password Reset Link to %s: %s\n", toEmail, resetLink)
		return nil
	}

	from := mail.NewEmail("MasterGo", s.fromEmail)
	to := mail.NewEmail("", toEmail)

	subject := "Reset Your Password - MasterGo"
	plainTextContent := fmt.Sprintf("Click here to reset your password: %s. Link expires in 1 hour.", resetLink)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Password Reset Request</h2>
				<p>You requested a password reset. Click the link below to reset your password:</p>
				<a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
					Reset Password
				</a>
				<p>This link will expire in 1 hour.</p>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
		</html>
	`, resetLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}

// Send2FADisabledEmail notifies user that 2FA was disabled
func (s *EmailService) Send2FADisabledEmail(toEmail string) error {
	if s.apiKey == "" {
		log.Printf("2FA Disabled Notification to %s\n", toEmail)
		return nil
	}

	from := mail.NewEmail("MasterGo", s.fromEmail)
	to := mail.NewEmail("", toEmail)

	subject := "Two-Factor Authentication Disabled - MasterGo"
	plainTextContent := "Two-factor authentication has been disabled on your account."
	htmlContent := `
		<html>
			<body>
				<h2>Security Alert</h2>
				<p>Two-factor authentication has been disabled on your account.</p>
				<p>If you didn't do this, please secure your account immediately.</p>
			</body>
		</html>
	`

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}

// SendAccountRecoveryEmail sends account recovery options
func (s *EmailService) SendAccountRecoveryEmail(toEmail, recoveryCode string) error {
	if s.apiKey == "" {
		log.Printf("Account Recovery Code to %s: %s\n", toEmail, recoveryCode)
		return nil
	}

	from := mail.NewEmail("MasterGo", s.fromEmail)
	to := mail.NewEmail("", toEmail)

	subject := "Account Recovery Code - MasterGo"
	plainTextContent := fmt.Sprintf("Your account recovery code is: %s", recoveryCode)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Account Recovery</h2>
				<p>Your recovery code is:</p>
				<h1 style="color: #007bff;">%s</h1>
				<p>Use this code to recover your account access.</p>
			</body>
		</html>
	`, recoveryCode)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}

// SendSecurityAlertEmail notifies user of suspicious activity
func (s *EmailService) SendSecurityAlertEmail(toEmail, action, ipAddress string) error {
	if s.apiKey == "" {
		log.Printf("Security Alert to %s: %s from %s\n", toEmail, action, ipAddress)
		return nil
	}

	from := mail.NewEmail("MasterGo", s.fromEmail)
	to := mail.NewEmail("", toEmail)

	subject := "Security Alert - Unusual Activity Detected"
	plainTextContent := fmt.Sprintf("Unusual activity detected on your account: %s from IP %s", action, ipAddress)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Security Alert</h2>
				<p><strong>Unusual activity detected on your account:</strong></p>
				<p>Action: %s</p>
				<p>IP Address: %s</p>
				<p>If this wasn't you, please change your password immediately.</p>
			</body>
		</html>
	`, action, ipAddress)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}
