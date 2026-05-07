package service

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	// SendGrid config
	apiKey    string
	fromEmail string

	// SMTP config for real email sending (Nodemailer-like)
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	useSMTP      bool
}

func NewEmailService(apiKey, fromEmail string) *EmailService {
	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}
}

// NewEmailServiceWithSMTP creates email service with SMTP configuration (for real email sending)
func NewEmailServiceWithSMTP(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail string) *EmailService {
	return &EmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
		useSMTP:      true,
	}
}

// sendViaSMTP sends email using SMTP (real email sending)
func (s *EmailService) sendViaSMTP(toEmail, subject, htmlContent, plainTextContent string) error {
	if s.smtpHost == "" || s.smtpPort == "" {
		return fmt.Errorf("SMTP configuration not set")
	}

	// Create message with both plain text and HTML
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=boundary123\r\n\r\n--boundary123\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s\r\n\r\n--boundary123\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n\r\n--boundary123--\r\n",
		s.fromEmail,
		toEmail,
		subject,
		plainTextContent,
		htmlContent,
	)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	addr := s.smtpHost + ":" + s.smtpPort

	err := smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("SMTP Error: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s via SMTP", toEmail)
	return nil
}

// SendOTPEmail sends OTP verification email via SMTP or SendGrid
func (s *EmailService) SendOTPEmail(toEmail, otp string) error {
	subject := "Verify Your Email - MasterGo"
	plainTextContent := fmt.Sprintf("Your verification code is: %s. It will expire in 5 minutes.", otp)
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

	// Try SMTP first (real email sending)
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil // Success
		}
		log.Printf("SMTP send failed, error: %v", err)
		// Fall through to SendGrid or console fallback
	}

	// Fall back to SendGrid if configured
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback for development
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | OTP: %s", toEmail, subject, otp)
	return nil
}

// SendPasswordResetEmail sends password reset link via SMTP or SendGrid
func (s *EmailService) SendPasswordResetEmail(toEmail, resetLink string) error {
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

	// Try SMTP first
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil
		}
		log.Printf("SMTP send failed, error: %v", err)
	}

	// Fall back to SendGrid
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | Link: %s", toEmail, subject, resetLink)
	return nil
}

// SendAccountRecoveryEmail sends account recovery options via SMTP or SendGrid
func (s *EmailService) SendAccountRecoveryEmail(toEmail, recoveryCode string) error {
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

	// Try SMTP first
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil
		}
		log.Printf("SMTP send failed, error: %v", err)
	}

	// Fall back to SendGrid
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | Code: %s", toEmail, subject, recoveryCode)
	return nil
}

// SendSecurityAlertEmail notifies user of suspicious activity via SMTP or SendGrid
func (s *EmailService) SendSecurityAlertEmail(toEmail, action, ipAddress string) error {
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

	// Try SMTP first
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil
		}
		log.Printf("SMTP send failed, error: %v", err)
	}

	// Fall back to SendGrid
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | Action: %s | IP: %s", toEmail, subject, action, ipAddress)
	return nil
}

// SendWelcomeEmail sends welcome email after successful signup
func (s *EmailService) SendWelcomeEmail(toEmail, userName string) error {
	subject := "Welcome to MasterGo!"
	plainTextContent := fmt.Sprintf("Welcome %s! Your account has been successfully created. You can now access all features of MasterGo.", userName)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>Welcome to MasterGo, %s!</h2>
				<p>Your account has been successfully created and verified.</p>
				<p>You can now:</p>
				<ul>
					<li>Generate and manage API keys</li>
					<li>Access your wallet and manage transactions</li>
					<li>Enable two-factor authentication for enhanced security</li>
					<li>Track your API usage and analytics</li>
				</ul>
				<p>If you have any questions, please contact our support team.</p>
				
			</body>
		</html>
	`, userName)

	// Try SMTP first
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil
		}
		log.Printf("SMTP send failed, error: %v", err)
	}

	// Fall back to SendGrid
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | User: %s", toEmail, subject, userName)
	return nil
}

// SendLoginNotificationEmail sends notification email on successful login
func (s *EmailService) SendLoginNotificationEmail(toEmail, userName, ipAddress string) error {
	subject := "New Login to Your MasterGo Account"
	plainTextContent := fmt.Sprintf("Hello %s, you successfully logged in to your MasterGo account from IP: %s", userName, ipAddress)
	htmlContent := fmt.Sprintf(`
		<html>
			<body>
				<h2>New Login Detected</h2>
				<p>Hello %s,</p>
				<p>You successfully logged in to your MasterGo account.</p>
				<p><strong>Login Details:</strong></p>
				<ul>
					<li>IP Address: %s</li>
					<li>Time: Just now</li>
				</ul>
				<p>If this wasn't you, please change your password immediately and contact our support.</p>
			</body>
		</html>
	`, userName, ipAddress)

	// Try SMTP first
	if s.useSMTP {
		err := s.sendViaSMTP(toEmail, subject, htmlContent, plainTextContent)
		if err == nil {
			return nil
		}
		log.Printf("SMTP send failed, error: %v", err)
	}

	// Fall back to SendGrid
	if s.apiKey != "" {
		from := mail.NewEmail("MasterGo", s.fromEmail)
		to := mail.NewEmail("", toEmail)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(s.apiKey)
		_, err := client.Send(message)
		if err == nil {
			log.Printf("Email sent successfully to %s via SendGrid", toEmail)
			return nil
		}
		log.Printf("SendGrid send failed: %v", err)
	}

	// Console fallback
	log.Printf("📧 [EMAIL FALLBACK - CONSOLE] To: %s | Subject: %s | User: %s | IP: %s", toEmail, subject, userName, ipAddress)
	return nil
}
