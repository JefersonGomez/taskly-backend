package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendNotificationEmail(to, nombre, equipoNombre, rol, invitedBy string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY no configurada")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://taskly-frontend-ruby.vercel.app"
	}

	client := resend.NewClient(apiKey)

	subject := fmt.Sprintf("🎉 Has sido agregado al equipo '%s' en Taskly", equipoNombre)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><style>
  body { font-family: -apple-system, sans-serif; line-height: 1.6; color: #172B4D; margin: 0; }
  .container { max-width: 600px; margin: 0 auto; padding: 24px; }
  .header { background: #0052CC; padding: 20px; text-align: center; border-radius: 12px 12px 0 0; }
  .header h1 { color: white; margin: 0; }
  .content { background: #fff; padding: 24px; border: 1px solid #DFE1E6; border-top: none; border-radius: 0 0 12px 12px; }
  .badge { display: inline-block; padding: 4px 12px; border-radius: 999px; font-size: 12px; font-weight: 600; background: #E9F2FF; color: #0052CC; }
  .btn { display: inline-block; background: #0052CC; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; margin: 16px 0; }
</style></head>
<body>
  <div class="container">
    <div class="header"><h1>🚀 Taskly</h1></div>
    <div class="content">
      <h2>¡Hola, %s! 👋</h2>
      <p><strong>%s</strong> te agregó al equipo <strong>"%s"</strong> como:</p>
      <p><span class="badge">%s</span></p>
      <a href="%s" class="btn">Ir a Taskly</a>
    </div>
  </div>
</body>
</html>`, nombre, invitedBy, equipoNombre, rol, appURL)

	params := &resend.SendEmailRequest{
		From:    "Taskly <onboarding@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	resp, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("❌ Resend error: %v", err)
		return err
	}

	log.Printf("✅ Email enviado a %s (ID: %s)", to, resp.Id)
	return nil
}
