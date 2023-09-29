package mail

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/teamhanko/hanko/backend/config"
	"github.com/teamhanko/hanko/backend/persistence/models"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	renderer    *Renderer
	emailConfig config.Email
	mailer      Mailer
}

type SendPasscodeEmailData struct {
	Code        string
	ServiceName string
	TTL         string
}

func NewNotificationService(cfg *config.Config, mailer Mailer) (*NotificationService, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to create new renderer: %w", err)
	}
	return &NotificationService{
		renderer:    renderer,
		emailConfig: cfg.Passcode.Email,
		mailer:      mailer,
	}, nil
}

func (m *NotificationService) SendPasscodeEmail(c echo.Context, toEmail *models.Email, data SendPasscodeEmailData) error {
	lang := c.Request().Header.Get("Accept-Language")
	str, err := m.renderer.Render("loginTextMail", lang, structs.Map(data))
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	message := gomail.NewMessage()
	message.SetAddressHeader("To", toEmail.Address, "")
	message.SetAddressHeader("From", m.emailConfig.FromAddress, m.emailConfig.FromName)

	message.SetHeader("Subject", m.renderer.Translate(lang, "email_subject_login", structs.Map(data)))

	message.SetBody("text/plain", str)

	err = m.mailer.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send passcode: %w", err)
	}
	return nil
}
