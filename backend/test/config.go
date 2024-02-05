package test

import "github.com/teamhanko/hanko/backend/config"

var DefaultConfig = config.Config{
	Webauthn: config.WebauthnSettings{
		RelyingParty: config.RelyingParty{
			Id:          "localhost",
			DisplayName: "Test Relying Party",
			Icon:        "",
			Origins:     []string{"http://localhost:8080", "http://localhost:8888"},
		},
		Timeout:          60000,
		UserVerification: "preferred",
	},
	Secrets: config.Secrets{
		Keys: []string{"abcdefghijklmnop"},
	},
	Smtp: config.SMTP{
		Host: "localhost",
		Port: "2500",
	},
	Passcode: config.Passcode{
		Email: config.Email{
			FromAddress: "test@hanko.io",
			FromName:    "Hanko Test",
		},
		TTL: 300,
	},
	Session: config.Session{
		Lifespan: "1h",
		Cookie: config.Cookie{
			SameSite: "none",
		},
	},
	Service: config.Service{
		Name: "Test",
	},
	Account: config.Account{
		AllowSignup:   true,
		AllowDeletion: false,
	},
	SecurityNotifications: config.SecurityNotifications{
		Notifications: config.SecurityNotificationTypes{
			PasswordUpdate: config.SecurityNotificationConfiguration{
				Enabled: false,
			},
			PrimaryEmailUpdate: config.SecurityNotificationConfiguration{
				Enabled: false,
			},
			EmailCreate: config.SecurityNotificationConfiguration{
				Enabled: false,
			},
			PasskeyCreate: config.SecurityNotificationConfiguration{
				Enabled: false,
			},
		},
		FromEmail: config.Email{
			FromAddress: "security@hanko.io",
			FromName:    "Hanko Security",
		},
		NotifyAddress: "help@hanko.io",
	},
}
