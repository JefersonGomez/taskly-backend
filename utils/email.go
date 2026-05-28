package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// SendNotificationEmail envía email usando SSL (puerto 465) o STARTTLS (puerto 587)
func SendNotificationEmail(to, nombre, equipoNombre, rol, invitedBy string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://taskly-frontend-ruby.vercel.app"
	}

	subject := fmt.Sprintf("🎉 Has sido agregado al equipo '%s' en Taskly", equipoNombre)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><style>
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #172B4D; }
  .container { max-width: 600px; margin: 0 auto; padding: 24px; }
  .header { background: linear-gradient(135deg, #0052CC, #00B8D9); padding: 20px; border-radius: 12px 12px 0 0; text-align: center; }
  .header h1 { color: white; margin: 0; font-size: 24px; }
  .content { background: #fff; padding: 24px; border: 1px solid #DFE1E6; border-top: none; border-radius: 0 0 12px 12px; }
  .badge { display: inline-block; padding: 4px 12px; border-radius: 999px; font-size: 12px; font-weight: 600; }
  .badge-admin { background: #E9F2FF; color: #0052CC; }
  .badge-member { background: #F4F5F7; color: #5E6C84; }
  .btn { display: inline-block; background: #0052CC; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; font-weight: 600; margin: 16px 0; }
  .footer { text-align: center; padding: 16px; color: #5E6C84; font-size: 12px; }
</style></head>
<body>
  <div class="container">
    <div class="header"><h1>🚀 Taskly</h1></div>
    <div class="content">
      <h2>¡Hola, %s! 👋</h2>
      <p><strong>%s</strong> te ha agregado al equipo <strong>"%s"</strong> con el rol de:</p>
      <p><span class="badge %s">%s</span></p>
      <p>Ahora puedes colaborar en tareas, actualizar estados y más.</p>
      <a href="%s" class="btn">Ir a Taskly</a>
    </div>
    <div class="footer"><p>© 2026 Taskly. Mensaje automático.</p></div>
  </div>
</body>
</html>`,
		nombre, invitedBy, equipoNombre,
		map[string]string{"admin": "badge-admin", "miembro": "badge-member"}[rol],
		rol, appURL)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body)

	addr := host + ":" + port

	// 🔹 Si usas puerto 465 → TLS implícito
	if port == "465" {
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("tls.Dial: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("smtp.NewClient: %w", err)
		}
		defer client.Close()

		if err = client.Auth(smtp.PlainAuth("", from, password, host)); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
		if err = client.Mail(from); err != nil {
			return fmt.Errorf("mail from: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("rcpt to: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("data: %w", err)
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("close data: %w", err)
		}
		return client.Quit()
	}

	// 🔹 Si usas puerto 587 → STARTTLS (comportamiento original)
	auth := smtp.PlainAuth("", from, password, host)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
