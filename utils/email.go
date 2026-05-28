package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendNotificationEmail envía una notificación HTML al usuario agregado
func SendNotificationEmail(to, nombre, equipoNombre, rol, invitedBy string) error {
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://taskly-frontend-ruby.vercel.app"
	}

	// Asunto y cuerpo HTML
	subject := fmt.Sprintf("🎉 Has sido agregado al equipo '%s' en Taskly", equipoNombre)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
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
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🚀 Taskly</h1>
    </div>
    <div class="content">
      <h2>¡Hola, %s! 👋</h2>
      <p><strong>%s</strong> te ha agregado al equipo <strong>"%s"</strong> con el rol de:</p>
      <p><span class="badge %s">%s</span></p>
      <p>Ahora puedes:</p>
      <ul>
        <li>✅ Ver y gestionar las tareas del equipo</li>
        <li>✅ Colaborar con otros miembros</li>
        <li>✅ Actualizar el estado de tus asignaciones</li>
      </ul>
      <a href="%s" class="btn">Ir a Taskly</a>
      <p style="font-size: 13px; color: #5E6C84; margin-top: 24px;">
        Si no esperabas esta invitación, puedes ignorar este correo o contactar al propietario del equipo.
      </p>
    </div>
    <div class="footer">
      <p>© 2026 Taskly. Este es un mensaje automático, por favor no responder.</p>
    </div>
  </div>
</body>
</html>
`,
		nombre,
		invitedBy,
		equipoNombre,
		map[string]string{"admin": "badge-admin", "miembro": "badge-member"}[rol],
		rol,
		appURL,
	)

	// Construir mensaje MIME
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body,
	)

	// Enviar
	err := smtp.SendMail(
		host+":"+port,
		smtp.PlainAuth("", from, password, host),
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("❌ Error enviando email a %s: %v", to, err)
		return err
	}

	log.Printf("✅ Email enviado a %s para equipo '%s'", to, equipoNombre)
	return nil
}
