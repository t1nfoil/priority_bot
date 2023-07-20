package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
)

// This was code was (mostly) modified from the slash example
// https://github.com/slack-go/slack/blob/master/examples/slash/slash.go

func verifySigningSecret(r *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, slackConfiguration.SigningSecret)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	verifier.Write(body)

	if err = verifier.Ensure(); err != nil {
		return err
	}

	return nil
}

func handleSlackSlashCommand(w http.ResponseWriter, r *http.Request) {

	err := verifySigningSecret(r)

	if err != nil {
		botLog(r, "error verifying signing secrets\n"+err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	botLog(s, "SlashCommandParse(s)")

	if err != nil {
		botLog(r, "error parsing slack command\n"+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	api := slack.New(slackConfiguration.ClientSecret)

	requestUser, err := api.GetUserInfo(s.UserID)

	if err != nil {
		botLog(api, "error getting user profile for requesting user ["+s.UserID+"]\n"+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch s.Command {

	case "/priority_bot":
		switch strings.ToLower(s.Text) {

		case "show-jira-fields":

			if botConfiguration.Debug == "true" {
				client, err := jiraInstance.NewAuthorizedClient()

				if err != nil {
					botLog(err, "jiraInstance.NewAuthorizedClient")
					return
				}

				fields, err := getJiraListOfIssueFields(client)

				if err != nil {
					botLog(err, "getJiraListOfIssueFields")
					return
				}

				botLog(fields, "Jira Fields")
			}

		case "check-token":

			if botConfiguration.Debug == "true" {

				// expire the token
				jiraInstance.Auth.token.Expiry = time.Now()
				jiraClient, err := jiraInstance.NewAuthorizedClient()

				//sendSlackMessage(requestUser.Email, "Error (creating Authorized Client): \n"+spew.Sdump(err))
				if err != nil {
					botLog(err, "jiraInstance.NewAuthorizedClient")
					botLog(spew.Sdump(jiraInstance.Auth.token), "Token Generation (failed)")
					return
				} else {
					botLog(string("Received Token"), "Token Generation (success)")
				}
				_, err = getJiraUserByEmail(jiraClient, requestUser.Profile.Email)
				//sendSlackMessage(requestUser.Email, "Error (attempting Jira Request): \n"+spew.Sdump(err))
				if err != nil {
					botLog(err, "getUserByEmail (failed)")
					botLog(requestUser.Profile.Email, "User Details")
					botLog(spew.Sdump(jiraClient), "Client Details")
					return
				} else {
					botLog(string("Token Functional"), "Token Test (success)")
				}
			} else {
				sendSlackMessage(requestUser.Profile.Email, "run `/priority_bot enable-debug` to use `refresh-token`")
			}
		case "enable-debug":
			for _, v := range botConfiguration.SlackControlUsers {
				if v == requestUser.Profile.Email {
					botConfiguration.Debug = "true"
					err := sendSlackMessage(requestUser.Profile.Email, "debug enabled")
					if err != nil {
						botLog(requestUser, err.Error())
					}
				}
			}
			w.WriteHeader(http.StatusOK)

		case "disable-debug":
			for _, v := range botConfiguration.SlackControlUsers {
				if v == requestUser.Profile.Email {
					err := sendSlackMessage(requestUser.Profile.Email, "debug disabled")
					if err != nil {
						botLog(requestUser, err.Error())
					}
					botConfiguration.Debug = "false"
				}
			}
			w.WriteHeader(http.StatusOK)

		case "help":
			botLog(nil, strings.ToLower(s.Text))
			// display help commands to user -- finalized later.

		case "auth-to-jira":
			botLog(nil, strings.ToLower(s.Text))

			err = jiraInstance.getAuthorization()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				botLog(nil, err.Error())
				return
			}

			w.WriteHeader(http.StatusOK)
			return

		default:

			var newTicket ticketRequest

			newTicket.toggleRollStatus = OPS_ROLL_INITIAL
			newTicket.toggleDeployFlavour = OPS_JIRA_DEPLOY_FLAVOR_INITIAL

			modalRequest := generateModal(newTicket.toggleRollStatus, newTicket.toggleDeployFlavour, newTicket.toggleProject)
			resp, err := api.OpenView(s.TriggerID, modalRequest)

			if err != nil {
				botLog(err, "error opening view")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			newTicket.SetViewID(resp.View.ID)
			newTicket.ticket.SetReporter(requestUser.Profile.Email)
			opsRequest.Add(newTicket)
			w.WriteHeader(http.StatusOK)
		}

	default:

		botLog(nil, "error bad slash command")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func handleSlackModal(w http.ResponseWriter, r *http.Request) {

	err := verifySigningSecret(r)

	if err != nil {
		botLog(err.Error(), "error verifying signing secret")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var slackCallBack slack.InteractionCallback

	err = json.Unmarshal([]byte(r.FormValue("payload")), &slackCallBack)

	if err != nil {
		botLog(err.Error(), "error unmarshaling callback payload")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch slackCallBack.Type {

	case slack.InteractionTypeBlockActions:
		err = handleBlockAction(slackCallBack)

	case slack.InteractionTypeViewClosed:
		err = handleViewClose(slackCallBack)

	case slack.InteractionTypeViewSubmission:
		err = handleViewSubmission(slackCallBack)
	}

	if err != nil {
		botLog(err.Error(), "handler error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleBlockAction(slackCallBack slack.InteractionCallback) error {

	ticket, err := opsRequest.GetTicketRequestByViewID(slackCallBack.View.ID)

	if err != nil {
		botLog(nil, err.Error())
		return err
	}

	blockActions := slackCallBack.ActionCallback.BlockActions

	for _, actionType := range blockActions {
		switch actionType.ActionID {

		case "a_id_rollstatus":

			if actionType.SelectedOption.Value == string(OPS_ROLL_BACK) {
				ticket.toggleRollStatus = OPS_ROLL_BACK
			}

			if actionType.SelectedOption.Value == string(OPS_ROLL_OUT) {
				ticket.toggleRollStatus = OPS_ROLL_OUT
			}

			if actionType.SelectedOption.Value == string(OPS_ROLL_INITIAL) {
				ticket.toggleRollStatus = OPS_ROLL_INITIAL
			}

		case "a_id_jira_project":

			if actionType.SelectedOption.Value == string(OPS_JIRA_PROJECT_OPS) {
				ticket.toggleProject = OPS_JIRA_PROJECT_OPS
			}

			if actionType.SelectedOption.Value == string(OPS_JIRA_PROJECT_OPSWORK) {
				ticket.toggleProject = OPS_JIRA_PROJECT_OPSWORK
			}

		case "a_id_jira_ticket":

			if actionType.SelectedOption.Value == string(OPS_JIRA_TICKET_TYPE_DEPLOY) {
				ticket.toggleDeployFlavour = OPS_JIRA_DEPLOY_FLAVOR_SHOW
			}

			if actionType.SelectedOption.Value == string(OPS_JIRA_TICKET_TYPE_MAINTENANCE) {
				ticket.toggleDeployFlavour = OPS_JIRA_DEPLOY_FLAVOR_INITIAL
			}

			if actionType.SelectedOption.Value == string(OPS_JIRA_TICKET_TYPE_TASK) {
				ticket.toggleDeployFlavour = OPS_JIRA_DEPLOY_FLAVOR_INITIAL
			}

			if actionType.SelectedOption.Value == string(OPS_JIRA_TICKET_TYPE_EPIC) {
				ticket.toggleDeployFlavour = OPS_JIRA_DEPLOY_FLAVOR_INITIAL
			}

		}
	}

	opsRequest.SetTicketRequestByViewID(ticket.viewID, ticket)

	var modalRequest slack.ModalViewRequest
	api := slack.New(slackConfiguration.ClientSecret)

	modalRequest = generateModal(ticket.toggleRollStatus, ticket.toggleDeployFlavour, ticket.toggleProject)
	_, err = api.UpdateView(modalRequest, slackCallBack.View.ExternalID, slackCallBack.Hash, slackCallBack.View.ID)

	if err != nil {
		botLog(slackCallBack, "error updating view")
		return err
	}

	return nil
}

func handleViewClose(slackCallBack slack.InteractionCallback) error {
	err := opsRequest.RemoveTicketRequestByViewID(slackCallBack.View.ID)

	if err != nil {
		return err
	}

	return nil
}

func handleViewSubmission(slackCallBack slack.InteractionCallback) error {

	submission, err := opsRequest.GetTicketRequestByViewID(slackCallBack.View.ID)

	if err != nil {
		return (err)
	}

	err = submission.ticket.PopulateTicket(&slackCallBack)

	if err != nil {
		return (err)
	}

	jiraClient, err := jiraInstance.NewAuthorizedClient()

	if err != nil {
		return (err)
	}

	createdIssue, err := createIssue(jiraClient, submission)

	if err != nil {

		// notify the user something went terribly wrong.

		err = sendSlackMessage(submission.ticket.Reporter, "This is embarrasing, but something it broken. We do have an error message that might assist in troubleshooting this: "+err.Error())

		if err != nil {
			botLog(err, "error sending slack message about jira issue failure.. things are really borky now")
		}

		return (err)
	}

	err = doSlackNotifications(submission, createdIssue)

	if err != nil {
		return err
	}

	err = opsRequest.RemoveTicketRequestByViewID(submission.viewID)

	if err != nil {
		botLog(submission, "error - the jira issue was created successfuly, but removing it from the queue failed")
		return err
	}

	return nil

}
