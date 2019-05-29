package slacklib

import (
	"golang.org/x/oauth2"

	oslack "golang.org/x/oauth2/slack"
)

type Conf struct {
	// Slack App Config
	AppID             string `json:"app_id"`
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	SigningSecret     string `json:"signing_secret"`
	VerificationToken string `json:"verification_token"` // deprecated in favour of SigningSecret
	BotToken          string `json:"bot_token,omitempty"`

	// OAuth config
	OAuthInstallConfig oauth2.Config
	OAuthLoginConfig   oauth2.Config
}

func (conf *Conf) Init() {
	conf.OAuthInstallConfig.ClientID = conf.ClientID
	conf.OAuthInstallConfig.ClientSecret = conf.ClientSecret
	conf.OAuthInstallConfig.Endpoint = oslack.Endpoint
	// .RedirectURL is set by user
	// .Scopes is set by user

	conf.OAuthLoginConfig.ClientID = conf.ClientID
	conf.OAuthLoginConfig.ClientSecret = conf.ClientSecret
	conf.OAuthLoginConfig.Endpoint = oslack.Endpoint
	// .RedirectURL is set by user
	// .Scopes is set by user
}
