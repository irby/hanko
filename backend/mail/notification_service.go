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

type SendPasscodeEmailData struct {
	Code        string
	ServiceName string
	TTL         string
}

func (m *NotificationService) SendPasscodeEmail(c echo.Context, toEmail *models.Email, data SendPasscodeEmailData) error {
	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		SubjectTemplate: "email_subject_login",
		BodyTemplate:    "loginTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, data)
}

type EmailProps struct {
	ToEmail         string
	SubjectTemplate string
	BodyTemplate    string
}

func (m *NotificationService) sendEmail(c echo.Context, props EmailProps, data interface{}) error {
	lang := c.Request().Header.Get("Accept-Language")
	body, err := m.renderer.Render(props.BodyTemplate, lang, structs.Map(data))
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	message := gomail.NewMessage()
	message.SetAddressHeader("To", props.ToEmail, "")
	message.SetAddressHeader("From", m.emailConfig.FromAddress, m.emailConfig.FromName)

	message.SetHeader("Subject", m.renderer.Translate(lang, props.SubjectTemplate, structs.Map(data)))

	message.SetBody("text/plain", body)

	err = m.mailer.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send passcode: %w", err)
	}
	return nil
}
