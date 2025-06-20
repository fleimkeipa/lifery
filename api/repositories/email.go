package repositories

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fleimkeipa/lifery/pkg"
	"gopkg.in/gomail.v2"
)

type EmailRepository struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailRepository() *EmailRepository {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := 587
	if port := os.Getenv("SMTP_PORT"); port != "" {
		if p, err := fmt.Sscanf(port, "%d", &smtpPort); err != nil || p != 1 {
			smtpPort = 587
		}
	}

	email := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	dialer := gomail.NewDialer(smtpHost, smtpPort, email, password)

	return &EmailRepository{
		dialer: dialer,
		from:   email,
	}
}

func (es *EmailRepository) SendPasswordResetEmail(to, username, resetToken string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:8081"
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, resetToken)

	subject := "Lifery - Şifre Sıfırlama"

	htmlBody := fmt.Sprintf(htmlBody, username, resetLink, resetLink)

	textBody := fmt.Sprintf(textBody, username, resetLink)

	m := gomail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	if err := es.dialer.DialAndSend(m); err != nil {
		return pkg.NewError(err, "failed to send password reset email", http.StatusInternalServerError)
	}

	return nil
}

var htmlBody = `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Şifre Sıfırlama</title>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
				.content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
				.button { display: inline-block; background-color: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
				.footer { text-align: center; margin-top: 30px; color: #666; font-size: 14px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Lifery</h1>
					<p>Şifre Sıfırlama İsteği</p>
				</div>
				<div class="content">
					<h2>Merhaba %s,</h2>
					<p>Lifery hesabınız için şifre sıfırlama isteği aldık.</p>
					<p>Şifrenizi sıfırlamak için aşağıdaki butona tıklayın:</p>
					
					<div style="text-align: center;">
						<a href="%s" class="button">Şifremi Sıfırla</a>
					</div>
					
					<p>Eğer bu isteği siz yapmadıysanız, bu emaili görmezden gelebilirsiniz.</p>
					<p>Bu link 24 saat boyunca geçerlidir.</p>
					
					<p>Eğer buton çalışmıyorsa, aşağıdaki linki tarayıcınıza kopyalayabilirsiniz:</p>
					<p style="word-break: break-all; color: #4F46E5;">%s</p>
				</div>
				<div class="footer">
					<p>Bu email Lifery uygulaması tarafından gönderilmiştir.</p>
					<p>© 2025 Lifery. Tüm hakları saklıdır.</p>
				</div>
			</div>
		</body>
		</html>
	`

var textBody = `
Şifre Sıfırlama İsteği

Merhaba %s,

Lifery hesabınız için şifre sıfırlama isteği aldık.

Şifrenizi sıfırlamak için aşağıdaki linke tıklayın:
%s

Eğer bu isteği siz yapmadıysanız, bu emaili görmezden gelebilirsiniz.
Bu link 24 saat boyunca geçerlidir.

Lifery Ekibi
`
