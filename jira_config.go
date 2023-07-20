package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type jiraConfig struct {
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
}

// Loads the token from the (user specifiable but defaulted to).gitignored token/ directory, allowing us to ensure the token does not get
// copied to the repo. If you fork this, you will need to get a bot token and save it inside the
// ./token/jira_tokens.json file or load appropriate env variables in your app.
//
// jira_tokens.json format
//
// {
//     "client_id": "clientid",
//     "client_secret": "yoursupersecretbottoken",
// }
//
//

func (s *jiraConfig) loadJiraConfiguration(tokenFilePath string) error {

	tokenFile, err := os.Open(path.Clean(tokenFilePath))

	if err != nil {
		var isClientSecretFound bool

		s.ClientID, _ = os.LookupEnv("BOT_JIRA_CLIENT_ID")
		s.ClientSecret, isClientSecretFound = os.LookupEnv("BOT_JIRA_CLIENT_SECRET")

		if !isClientSecretFound {
			return (errors.New("unable to load client secret from " + tokenFilePath + " or environment variable BOT__JIRA_CLIENT_SECRET"))
		}

		if isClientSecretFound && s.ClientSecret == "" {
			return (errors.New("unable to load client secret from " + tokenFilePath + " and environment variable BOT_JIRA_CLIENT_SECRET is empty"))
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
