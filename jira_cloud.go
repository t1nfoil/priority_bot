package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type jiraAuth struct {
	token       *oauth2.Token
	config      *oauth2.Config
	stateString string
}

type jiraCloudConfig []struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	AvatarURL string   `json:"avatarUrl"`
}

type jiraCloud struct {
	ApiURL string
	Auth   jiraAuth
	Config jiraCloudConfig
	ctx    context.Context
}

func (j *jiraCloud) setContext() {
	j.ctx = context.Background()
}

func (j *jiraCloud) configJiraAuth() {

	stateString := strings.Replace(uuid.NewString(), "-", "", -1)
	redirectURL := botConfiguration.JiraRedirectUrl
	tokenURL := botConfiguration.JiraTokenUrl
	authURL := botConfiguration.JiraAuthUrl + "&state=" + stateString

	jiraEndPoint := oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}

	j.Auth.config = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     jiraConfiguration.ClientID,
		ClientSecret: jiraConfiguration.ClientSecret,
		Endpoint:     jiraEndPoint,
	}

	j.Auth.stateString = stateString

}

func (j *jiraCloud) getAuthorization() error {
	var reqBody bytes.Buffer
	request, err := http.NewRequest("GET", jiraInstance.Auth.config.Endpoint.AuthURL, &reqBody)

	if err != nil {
		botLog(request, "error creating authorization request", err.Error())
		return err
	}

	var client http.Client

	_, err = client.Do(request)

	if err != nil {
		botLog(err.Error(), "error receiving authorization response")
		return err
	}

	err = sendSlackMessage(botConfiguration.JiraOauthUser, "Please click on the following auth url to allow priority_bot to connect to Jira\n"+jiraInstance.Auth.config.Endpoint.AuthURL)

	if err != nil {
		botLog(err, "error sending auth url via slack")
		return err
	}

	return nil
}

func (j *jiraCloud) Token() (*oauth2.Token, error) {
	token := j.Auth.token
	conf := j.Auth.config

	if token.Expiry.Before(time.Now()) { // expired so let's update it
		src := conf.TokenSource(j.ctx, token)

		newToken, err := src.Token() // this actually goes and renews the tokens
		if err != nil {
			botLog(err, "error renewing token")
			return nil, errors.New("error renewing token")
		}

		if newToken.AccessToken != token.AccessToken {
			j.Auth.token = newToken
			botLog(nil, "setting new token")
			return newToken, nil
		}
	}

	return j.Auth.token, nil
}

func (j *jiraCloud) getCloudParameters() error {

	jiraClient := j.Auth.config.Client(j.ctx, j.Auth.token)
	response, err := jiraClient.Get("https://api.atlassian.com/oauth/token/accessible-resources")

	if err != nil {
		botLog(err, "error getting list of accessible resources")
		return (err)
	}

	defer response.Body.Close()
	stringBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		botLog(err, "error reading body response")
		return (err)
	}

	var jiraCloudParameters jiraCloudConfig
	err = json.Unmarshal(stringBytes, &jiraCloudParameters)

	if err != nil {
		botLog(err, "error unmarshalling cloud parameters from body response")
		return (err)
	}

	jiraInstance.Config = jiraCloudParameters

	jiraInstance.ApiURL = "https://api.atlassian.com/ex/jira/" + jiraCloudParameters[0].ID + "/"

	return nil
}

// makes authorized client
func (j *jiraCloud) NewAuthorizedClient() (*jira.Client, error) {

	sourceToken := oauth2.ReuseTokenSource(j.Auth.token, j)
	token, err := sourceToken.Token()

	if err != nil {
		botLog(err, "error obtaining token")
		return &jira.Client{}, err
	}

	jiraClient, err := jira.NewClient(j.Auth.config.Client(j.ctx, token), j.ApiURL)

	if err != nil {
		botLog(err, "error creating a new authorized jira client")
		return &jira.Client{}, err
	}

	return jiraClient, nil

}

func handleJiraCallback(w http.ResponseWriter, r *http.Request) {

	receivedStateString := r.FormValue("state")
	code := r.FormValue("code")

	if receivedStateString != jiraInstance.Auth.stateString {
		botLog(nil, "invalid state received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	jiraInstance.Auth.token, err = jiraInstance.Auth.config.Exchange(jiraInstance.ctx, code, oauth2.AccessTypeOffline)

	if err != nil {
		botLog(err, "token exchange failed")
		w.WriteHeader(http.StatusBadRequest)
		jiraInstance.Auth.token = nil
		return
	}

	if jiraInstance.Auth.token.RefreshToken == "" {
		botLog(nil, "handleJiraCallback has an empty refresh token")
	}

	// sending bytes will send status code 200 automatically
	w.Write([]byte("<html><head></head><body><h3>priority bot has been authorized under your account</h3>\n\n<code>01010100 01101000 01100001 01101110\n01101011 01110011 00100001 00100000\n00111010 00101001 00000000 00000000</code></body></html>"))

	err = jiraInstance.getCloudParameters()

	if err != nil {
		botLog(nil, err.Error())
		botLog(jiraInstance.Config, "jiraInstance.Config below")
		return
	}
}
