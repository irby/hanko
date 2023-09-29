package mail

import (
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/teamhanko/hanko/backend/persistence/models"
	"github.com/teamhanko/hanko/backend/test"
	"net/http"
	"net/http/httptest"
	"regexp"
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
	service := s.getService()
	testContext := s.generateTestContext()
	testEmail := models.Email{
		Address: "test.sendpasscode@example.com",
	}
	testData := SendPasscodeEmailData{
		Code:        "12345",
		ServiceName: "Test Service",
	}
	err := service.SendPasscodeEmail(testContext, &testEmail, testData)
	s.NoError(err)
	toAddress, subject := s.getLastEmailToAddressAndSubject()
	s.Equal(testEmail.Address, toAddress)
	s.Equal(s.translate(service, "email_subject_login", testData), subject)
}

func (s *notificationServiceSuite) getService() *NotificationService {
	mailer, err := NewMailer(test.DefaultConfig.Passcode.Smtp)
	s.Require().NoError(err)
	service, err := NewNotificationService(&test.DefaultConfig, mailer)
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
