package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type botConfig struct {
	SlackControlUsers []string `json:"control_users"`
	JiraOauthUser     string   `json:"jira_oauth_user"`
	JiraAuthUrl       string   `json:"jira_auth_url"`
	JiraTokenUrl      string   `json:"jira_token_url"`
	JiraRedirectUrl   string   `json:"jira_redirect_url"`
	Debug             string   `json:"debug"`
}

// Loads the configuration for bot parameters (the users who are allowed to execute admin functions on the bot), plus jira URLs.
// control_users are users who can run admin functions from the slack interface (such as the slack command /priority_boy auth-to-jira)
// jira_oauth_user is the slack person who will auth the bot for jira scopes + access token, this person must have appropriate permissions
// jira_auth_url, jira_token_url and jira_redirect_url are the urls from published app settings -- see README on how to set these.
//
// bot_config.json format
//
// {
//     "control_users": [
//         "user1@platform.sh",
//         "user2@platform.sh"
//     ]
//     "jira_oauth_user": "oauthuser@platform.sh"
//     "jira_auth_url": "https://link.to.authurl",
//     "jira_token_url": "https://auth.atlassian.com/oauth/token",
//     "jira_redirect_url": "https://link.to.redirecturl",
//     "debug":"false"
// }
//
//

func (b *botConfig) loadBotConfiguration(configFilePath string) error {

	configFile, err := os.Open(path.Clean(configFilePath))

	if err != nil {
		return err

	}

	configFileData, err := ioutil.ReadAll(configFile)

	if err != nil {
		return err
	}

	err = json.Unmarshal(configFileData, &b)

	if err != nil {
		return err
	}

	if len(b.SlackControlUsers) < 1 {
		return errors.New("error in ./config/bot_config.json: no control_users specified")
	}

	if b.JiraOauthUser == "" {
		return errors.New("error in ./config/bot_config.json: no jira_oauth_user specified")
	}

	if b.JiraAuthUrl == "" {
		return errors.New("error in ./config/bot_config.json: no jira_auth_url specified")
	}

	if b.JiraTokenUrl == "" {
		return errors.New("error in ./config/bot_config.json: no jira_token_url specified")
	}

	if b.JiraRedirectUrl == "" {
		return errors.New("error in ./config/bot_config.json: no jira_redirect_url specified")
	}

	return nil
}
