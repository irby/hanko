package mail

import (
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/teamhanko/hanko/backend/config"
	"github.com/teamhanko/hanko/backend/persistence/models"
	"github.com/teamhanko/hanko/backend/test"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestNewNotificationService(t *testing.T) {
	t.Parallel()
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
	service := s.getService(&test.DefaultConfig)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasscode@example.com",
	}
	testData := SendPasscodeEmailData{
		Code: "12345",
		TTL:  "5 minutes",
	}
	countBefore := s.getMessageCount()
	err := service.SendPasscodeEmail(testContext, &testEmail, testData)
	s.NoError(err)
	s.Equal(countBefore+1, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasswordUpdateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PasswordUpdate.Enabled = false
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasswordupdate@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendPasswordUpdateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasswordUpdateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PasswordUpdate.Enabled = true
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasswordupdate@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendPasswordUpdateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore+1, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPrimaryEmailUpdateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PrimaryEmailUpdate.Enabled = false
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendprimaryemailupdate@example.com",
	}
	data := SendPrimaryEmailUpdateEmailData{
		OldEmailAddress: "foo@hanko.io",
		NewEmailAddress: "bar@hanko.io",
	}
	countBefore := s.getMessageCount()
	err := service.SendPrimaryEmailUpdateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(countBefore, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPrimaryEmailUpdateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PrimaryEmailUpdate.Enabled = true
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendprimaryemailupdate@example.com",
	}
	data := SendPrimaryEmailUpdateEmailData{
		OldEmailAddress: "foo@hanko.io",
		NewEmailAddress: "bar@hanko.io",
	}
	countBefore := s.getMessageCount()
	err := service.SendPrimaryEmailUpdateEmail(testContext, &testEmail, data)
	s.NoError(err)
	s.Equal(countBefore+1, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendEmailCreateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.EmailCreate.Enabled = false
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendemailcreateemail@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendEmailCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendEmailCreateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.EmailCreate.Enabled = true
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendemailcreateemail@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendEmailCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore+1, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasskeyCreateEmail_Disabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PasskeyCreate.Enabled = false
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasskeycreateemail@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendPasskeyCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore, s.getMessageCount())
}

func (s *notificationServiceSuite) TestNotificationService_SendPasskeyCreateEmail_Enabled() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode")
	}
	cfg := &test.DefaultConfig
	cfg.SecurityNotifications.Notifications.PasskeyCreate.Enabled = true
	service := s.getService(cfg)
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasskeycreateemail@example.com",
	}
	countBefore := s.getMessageCount()
	err := service.SendPasskeyCreateEmail(testContext, &testEmail)
	s.NoError(err)
	s.Equal(countBefore+1, s.getMessageCount())
}

func (s *notificationServiceSuite) getService(config *config.Config) *NotificationService {
	mailer, err := NewMailer(test.DefaultConfig.Passcode.Smtp)
	s.Require().NoError(err)
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

func (s *notificationServiceSuite) getLastEmailToAddressAndSubject() (string, string) {
	messages := s.EmailServer.Messages()
	lastMessage := messages[len(messages)-1]

	subject := regexp.MustCompile("Subject: (.*)\r")
	subjectMatchTo := subject.FindStringSubmatch(lastMessage.MsgRequest())

	toAddress := regexp.MustCompile("To: (.*)\r")
	toAddressMatchTo := toAddress.FindStringSubmatch(lastMessage.MsgRequest())

	return toAddressMatchTo[1], subjectMatchTo[1]
}

func (s *notificationServiceSuite) translate(service *NotificationService, messageId string, data interface{}) string {
	return service.renderer.Translate(lang, messageId, structs.Map(data))
}

func (s *notificationServiceSuite) getMessageCount() int {
	return len(s.EmailServer.Messages())
}
