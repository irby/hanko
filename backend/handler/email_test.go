package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/teamhanko/hanko/backend/crypto/jwk"
	"github.com/teamhanko/hanko/backend/dto"
	"github.com/teamhanko/hanko/backend/session"
	"github.com/teamhanko/hanko/backend/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailSuite(t *testing.T) {
	t.Parallel()
	s := new(emailSuite)
	s.WithEmailServer = true
	suite.Run(t, s)
}

type emailSuite struct {
	test.Suite
}

func (s *emailSuite) TestEmailHandler_List() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/email")
	s.Require().NoError(err)

	e := NewPublicRouter(&test.DefaultConfig, s.Storage, nil)

	jwkManager, err := jwk.NewDefaultManager(test.DefaultConfig.Secrets.Keys, s.Storage.GetJwkPersister())
	s.Require().NoError(err)
	sessionManager, err := session.NewManager(jwkManager, test.DefaultConfig)
	s.Require().NoError(err)

	tests := []struct {
		name          string
		userId        uuid.UUID
		expectedCount int
	}{
		{
			name:          "should return all user assigned email addresses",
			userId:        uuid.FromStringOrNil("b5dd5267-b462-48be-b70d-bcd6f1bbe7a5"),
			expectedCount: 3,
		},
		{
			name:          "should return an empty list when the user has no email addresses assigned",
			userId:        uuid.FromStringOrNil("d41df4b7-c055-45e6-9faf-61aa92a4032e"),
			expectedCount: 0,
		},
	}

	for _, currentTest := range tests {
		s.Run(currentTest.name, func() {
			token, err := sessionManager.GenerateJWT(currentTest.userId)
			s.Require().NoError(err)
			cookie, err := sessionManager.GenerateCookie(token)
			s.Require().NoError(err)

			req := httptest.NewRequest(http.MethodGet, "/emails", nil)
			req.AddCookie(cookie)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			if s.Equal(http.StatusOK, rec.Code) {
				var emails []*dto.EmailResponse
				s.NoError(json.Unmarshal(rec.Body.Bytes(), &emails))
				s.Equal(currentTest.expectedCount, len(emails))
			}
		})
	}
}

func (s *emailSuite) TestEmailHandler_Create() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/email")
	s.Require().NoError(err)

	config := test.DefaultConfig
	config.Emails.MaxNumOfAddresses = 100
	config.SecurityNotifications.Notifications.EmailCreate.Enabled = true

	e := NewPublicRouter(&config, s.Storage, nil)

	jwkManager, err := jwk.NewDefaultManager(config.Secrets.Keys, s.Storage.GetJwkPersister())
	s.Require().NoError(err)
	sessionManager, err := session.NewManager(jwkManager, config)
	s.Require().NoError(err)

	userId := uuid.FromStringOrNil("b5dd5267-b462-48be-b70d-bcd6f1bbe7a5")

	token, err := sessionManager.GenerateJWT(userId)
	s.NoError(err)
	cookie, err := sessionManager.GenerateCookie(token)
	s.NoError(err)

	body := dto.EmailCreateRequest{
		Address: "test@example.com",
	}
	bodyJson, err := json.Marshal(body)
	s.NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/emails", bytes.NewReader(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()

	countBefore := len(s.EmailServer.Messages())
	e.ServeHTTP(rec, req)
	if s.Equal(http.StatusOK, rec.Code) {
		s.Equal(countBefore+1, len(s.EmailServer.Messages()))
	}
}

func (s *emailSuite) TestEmailHandler_SetPrimaryEmail() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/email")
	s.Require().NoError(err)

	e := NewPublicRouter(&test.DefaultConfig, s.Storage, nil)

	jwkManager, err := jwk.NewDefaultManager(test.DefaultConfig.Secrets.Keys, s.Storage.GetJwkPersister())
	s.Require().NoError(err)
	sessionManager, err := session.NewManager(jwkManager, test.DefaultConfig)
	s.Require().NoError(err)

	oldPrimaryEmailId := uuid.FromStringOrNil("51b7c175-ceb6-45ba-aae6-0092221c1b84")
	newPrimaryEmailId := uuid.FromStringOrNil("8bb4c8a7-a3e6-48bb-b54f-20e3b485ab33")
	userId := uuid.FromStringOrNil("b5dd5267-b462-48be-b70d-bcd6f1bbe7a5")

	token, err := sessionManager.GenerateJWT(userId)
	s.NoError(err)
	cookie, err := sessionManager.GenerateCookie(token)
	s.NoError(err)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/emails/%s/set_primary", newPrimaryEmailId), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	if s.Equal(http.StatusNoContent, rec.Code) {
		emails, err := s.Storage.GetEmailPersister().FindByUserId(userId)
		s.Require().NoError(err)

		s.Equal(3, len(emails))
		for _, email := range emails {
			if email.ID == newPrimaryEmailId {
				s.True(email.IsPrimary())
			} else if email.ID == oldPrimaryEmailId {
				s.False(email.IsPrimary())
			}
		}
	}
}

func (s *emailSuite) TestEmailHandler_Delete() {
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/email")
	s.Require().NoError(err)

	e := NewPublicRouter(&test.DefaultConfig, s.Storage, nil)
	userId := uuid.FromStringOrNil("b5dd5267-b462-48be-b70d-bcd6f1bbe7a5")

	jwkManager, err := jwk.NewDefaultManager(test.DefaultConfig.Secrets.Keys, s.Storage.GetJwkPersister())
	s.Require().NoError(err)
	sessionManager, err := session.NewManager(jwkManager, test.DefaultConfig)
	s.Require().NoError(err)

	token, err := sessionManager.GenerateJWT(userId)
	s.NoError(err)
	cookie, err := sessionManager.GenerateCookie(token)
	s.NoError(err)

	tests := []struct {
		name          string
		emailId       uuid.UUID
		responseCode  int
		expectedCount int
	}{
		{
			name:          "should delete the email address",
			emailId:       uuid.FromStringOrNil("8bb4c8a7-a3e6-48bb-b54f-20e3b485ab33"),
			responseCode:  http.StatusNoContent,
			expectedCount: 2,
		},
		{
			name:          "should not delete the primary email address",
			emailId:       uuid.FromStringOrNil("51b7c175-ceb6-45ba-aae6-0092221c1b84"),
			responseCode:  http.StatusConflict,
			expectedCount: 2,
		},
	}

	for _, currentTest := range tests {
		s.Run(currentTest.name, func() {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/emails/%s", currentTest.emailId), nil)
			req.AddCookie(cookie)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			if s.Equal(currentTest.responseCode, rec.Code) {
				emails, err := s.Storage.GetEmailPersister().FindByUserId(userId)
				s.Require().NoError(err)
				s.Equal(currentTest.expectedCount, len(emails))
			}
		})
	}

}
