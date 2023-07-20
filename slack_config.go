package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type slackConfig struct {
	AppId         string `json:"app_id"`
	ClientSecret  string `json:"client_secret"`
	ClientId      string `json:"client_id"`
	SigningSecret string `json:"signing_secret"`
}

// Loads the token from the (user specifiable but defaulted to).gitignored token/ directory, allowing us to ensure the token does not get
// copied to the repo. If you fork this, you will need to get a bot token and save it inside the
// ./token/slack_tokens.json file or load appropriate env variables in your app.
//
// slack_tokens.json format
//
// {
//     "app_id": "",
//     "client_secret": "xoxb-yoursupersecretbottoken",
//     "client_id": "clientid",
//     "signing_secret": "signingsecret"
// }
//
//

func (s *slackConfig) loadSlackConfiguration(tokenFilePath string) error {

	tokenFile, err := os.Open(path.Clean(tokenFilePath))

	if err != nil {
		var isClientSecretFound bool

		s.AppId, _ = os.LookupEnv("BOT_SLACK_APP_ID")
		s.ClientSecret, isClientSecretFound = os.LookupEnv("BOT_SLACK_CLIENT_SECRET")
		s.ClientId, _ = os.LookupEnv("BOT_SLACK_CLIENT_ID")
		s.SigningSecret, _ = os.LookupEnv("BOT_SLACK_SIGNING_SECRET")

		if !isClientSecretFound {
			return (errors.New("unable to load client secret from " + tokenFilePath + " or environment variable BOT_SLACK_CLIENT_SECRET"))
		}

		if isClientSecretFound && s.ClientSecret == "" {
			return (errors.New("unable to load client secret from " + tokenFilePath + " and environment variable BOT_SLACK_CLIENT_SECRET is empty"))
		}

		return nil
	}

	tokenFileData, err := ioutil.ReadAll(tokenFile)

	if err != nil {
		return err
	}

	err = json.Unmarshal(tokenFileData, &s)

	if err != nil {
		return err
	}

	return nil
}
