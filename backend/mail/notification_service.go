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
	renderer           *Renderer
	config             *config.Config
	emailConfig        config.Email
	notificationConfig config.SecurityNotifications
	mailer             Mailer
}

func NewNotificationService(cfg *config.Config, mailer Mailer) (*NotificationService, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to create new renderer: %w", err)
	}
	return &NotificationService{
		renderer:           renderer,
		config:             cfg,
		emailConfig:        cfg.SecurityNotifications.FromEmail,
		notificationConfig: cfg.SecurityNotifications,
		mailer:             mailer,
	}, nil
}

type SendPasscodeEmailData struct {
	Code string
	TTL  string
}

func (m *NotificationService) SendPasscodeEmail(c echo.Context, toEmail *models.Email, data SendPasscodeEmailData) error {
	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		NotifyAddress:   m.notificationConfig.NotifyAddress,
		SubjectTemplate: "email_subject_login",
		BodyTemplate:    "loginTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, m.generateSendEmailData(data))
}

func (m *NotificationService) SendPasswordUpdateEmail(c echo.Context, toEmail *models.Email) error {
	if !m.notificationConfig.Notifications.PasswordUpdate.Enabled {
		return nil
	}

	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		NotifyAddress:   m.notificationConfig.NotifyAddress,
		SubjectTemplate: "email_subject_password_update",
		BodyTemplate:    "passwordUpdateTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, m.generateSendEmailData(struct{}{}))
}

type SendPrimaryEmailUpdateEmailData struct {
	OldEmailAddress string
	NewEmailAddress string
}

func (m *NotificationService) SendPrimaryEmailUpdateEmail(c echo.Context, toEmail *models.Email, data SendPrimaryEmailUpdateEmailData) error {
	if !m.notificationConfig.Notifications.PrimaryEmailUpdate.Enabled {
		return nil
	}

	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		NotifyAddress:   m.notificationConfig.NotifyAddress,
		SubjectTemplate: "email_subject_primary_email_update",
		BodyTemplate:    "primaryEmailUpdateTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, m.generateSendEmailData(data))
}

func (m *NotificationService) SendEmailCreateEmail(c echo.Context, toEmail *models.Email) error {
	if !m.notificationConfig.Notifications.EmailCreate.Enabled {
		return nil
	}

	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		NotifyAddress:   m.notificationConfig.NotifyAddress,
		SubjectTemplate: "email_subject_email_create",
		BodyTemplate:    "createEmailTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, m.generateSendEmailData(struct{}{}))
}

func (m *NotificationService) SendPasskeyCreateEmail(c echo.Context, toEmail *models.Email) error {
	if !m.notificationConfig.Notifications.PasskeyCreate.Enabled {
		return nil
	}

	var passcodeEmailProps = EmailProps{
		ToEmail:         toEmail.Address,
		NotifyAddress:   m.notificationConfig.NotifyAddress,
		SubjectTemplate: "email_subject_passkey_create",
		BodyTemplate:    "passkeyCreateTextMail",
	}
	return m.sendEmail(c, passcodeEmailProps, m.generateSendEmailData(struct{}{}))
}

type EmailProps struct {
	ToEmail         string
	NotifyAddress   string
	SubjectTemplate string
	BodyTemplate    string
}

func (m *NotificationService) sendEmail(c echo.Context, props EmailProps, data map[string]interface{}) error {
	lang := c.Request().Header.Get("Accept-Language")
	body, err := m.renderer.Render(props.BodyTemplate, lang, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	message := gomail.NewMessage()
	message.SetAddressHeader("To", props.ToEmail, "")
	message.SetAddressHeader("From", m.emailConfig.FromAddress, m.emailConfig.FromName)

	message.SetHeader("Subject", m.renderer.Translate(lang, props.SubjectTemplate, data))

	message.SetBody("text/plain", body)

	err = m.mailer.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send passcode: %w", err)
	}
	return nil
}

func (m *NotificationService) generateSendEmailData(data interface{}) map[string]interface{} {
	result := structs.Map(data)
	result["NotifyAddress"] = m.notificationConfig.NotifyAddress
	result["ServiceName"] = m.config.Service.Name
	return result
}
