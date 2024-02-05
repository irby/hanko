package service

import (
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/teamhanko/hanko/backend/config"
	"github.com/teamhanko/hanko/backend/mail"
	"github.com/teamhanko/hanko/backend/persistence/models"
	"github.com/teamhanko/hanko/backend/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewNotificationService(t *testing.T) {
	s := new(notificationServiceSuite)
	s.WithEmailServer = true
	suite.Run(t, s)
}

type notificationServiceSuite struct {
	test.Suite
}

const lang = "en-us"

func (s *notificationServiceSuite) TestNotificationService_SendPasscodeEmail() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	service := s.getService(s.getTestConfig(), testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasscode@example.com",
	}
	testData := SendPasscodeEmailData{
		Code: "12345",
		TTL:  "5 minutes",
	}
	err := service.SendPasscodeEmail(testContext, &testEmail, testData)
	s.NoError(err)
	s.Equal(1, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasswordUpdateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PasswordUpdate.Enabled = false
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasswordupdate@example.com",
	}
	err := service.SendPasswordUpdateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(0, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasswordUpdateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PasswordUpdate.Enabled = true
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasswordupdate@example.com",
	}
	err := service.SendPasswordUpdateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(1, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPrimaryEmailUpdateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PrimaryEmailUpdate.Enabled = false
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendprimaryemailupdate@example.com",
	}
	data := SendPrimaryEmailUpdateEmailData{
		OldEmailAddress: "foo@hanko.io",
		NewEmailAddress: "bar@hanko.io",
	}
	err := service.SendPrimaryEmailUpdateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(0, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPrimaryEmailUpdateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PrimaryEmailUpdate.Enabled = true
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendprimaryemailupdate@example.com",
	}
	data := SendPrimaryEmailUpdateEmailData{
		OldEmailAddress: "foo@hanko.io",
		NewEmailAddress: "bar@hanko.io",
	}
	err := service.SendPrimaryEmailUpdateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(1, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendEmailCreateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.EmailCreate.Enabled = false
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendemailcreateemail@example.com",
	}
	data := SendEmailCreateEmailData{
		NewEmailAddress: "test@example.com",
	}
	err := service.SendEmailCreateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(0, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendEmailCreateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.EmailCreate.Enabled = true
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendemailcreateemail@example.com",
	}
	data := SendEmailCreateEmailData{
		NewEmailAddress: "test@example.com",
	}
	err := service.SendEmailCreateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(1, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasskeyCreateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PasskeyCreate.Enabled = false
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasskeycreateemail@example.com",
	}
	err := service.SendPasskeyCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(0, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasskeyCreateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	testMailer := s.getTestMailer()
	cfg := s.getTestConfig()
	cfg.SecurityNotifications.Notifications.PasskeyCreate.Enabled = true
	service := s.getService(cfg, testMailer.GetMailer())
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasskeycreateemail@example.com",
	}
	err := service.SendPasskeyCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(1, testMailer.GetEmailCount())
}

func (s *notificationServiceSuite) getService(config *config.Config, mailer mail.Mailer) *NotificationService {
	service, err := NewNotificationService(config, mailer)
	s.Require().NoError(err)
	s.NotNil(service)
	return service
}

func (s *notificationServiceSuite) generateTestContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set("Accept-Language", lang)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func (s *notificationServiceSuite) translate(service *NotificationService, messageId string, data interface{}) string {
	return service.renderer.Translate(lang, messageId, structs.Map(data))
}

func (s *notificationServiceSuite) getTestConfig() *config.Config {
	cfg := &test.DefaultConfig
	cfg.Passcode.Smtp.Host = s.EmailServer.SmtpHost
	cfg.Passcode.Smtp.Port = s.EmailServer.SmtpPort
	return cfg
}

func (s *notificationServiceSuite) getTestMailer() *test.TestMailer {
	mailer := test.NewTestMailer()
	return mailer
}
