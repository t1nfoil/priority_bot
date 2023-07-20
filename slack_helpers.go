package main

import (
	"github.com/andygrunwald/go-jira"
	"github.com/slack-go/slack"
)

func sendSlackMessage(userEmail string, message string) error {

	botLog(message, "Sending slack message to "+userEmail)
	api := slack.New(slackConfiguration.ClientSecret)

	user, err := api.GetUserByEmail(userEmail)

	if err != nil {
		botLog(user, "error finding slack user for email address")
		return err
	}

	_, _, err = api.PostMessage(user.ID, slack.MsgOptionText(message, false))

	if err != nil {
		botLog(nil, "error sending auth url via slack")
		return err
	}
	return nil
}

func postMessageToOps(submission ticketRequest, opsName, opsLink string) error {

	p := submission.ticket.PriorityLevel

	if p == OPS_PRIORITY_P3 || p == OPS_PRIORITY_P4 || p == OPS_PRIORITY_P5 {
		botLog(nil, "Not a P1 or P2 ticket, not submitting to #ops")
		return nil
	}

	headerText := slack.NewTextBlockObject("mrkdwn", ":warning: *Ops Priority Request* :warning:", false, false)

	message := ""
	priorityString := ""

	message += "*Request received from " + submission.ticket.Reporter + "*\n"
	if submission.ticket.PriorityLevel == OPS_PRIORITY_P1 {
		priorityString = string(OPS_PRIORITY_P1_DESC)
	} else {
		priorityString = string(OPS_PRIORITY_P2_DESC)
	}
	// Priority Level
	message += "> *Priority Level:* " + priorityString + "\n"

	// Rollout or Rollback
	if submission.ticket.RollStatus == OPS_ROLL_BACK {
		message += "> *Roll Type:* " + string(submission.ticket.RollStatus) + "\n"
		if submission.ticket.RollBackScriptLink == "" {
			message += "> *Rollback Script Tested:* " + string(submission.ticket.RollBackScriptTested) + "\n"
		} else {
			message += "> *Rollback Script Tested:* " + string(submission.ticket.RollBackScriptTested) + "  (" + "<" + string(submission.ticket.RollBackScriptLink) + "|Rollback Script Link>" + ")\n"
		}
	} else {
		message += "> *Roll Type:* " + string(submission.ticket.RollStatus) + "\n"
		if submission.ticket.RollBackScriptLink == "" {
			message += "> *Rollout Script Tested:* " + string(submission.ticket.RollOutScriptTested) + "\n"
		} else {
			message += "> *Rollout Script Tested:* " + string(submission.ticket.RollOutScriptTested) + "  (" + "<" + string(submission.ticket.RollOutScriptLink) + "|Rollout Script Link>" + ")\n"
		}
	}

	// Request Title
	message += "> *Request:* " + string(submission.ticket.RequestTitle) + "\n"

	// Regions
	message += "> *Region(s):*\n"

	for _, v := range submission.ticket.Regions {
		message += "> - " + v + "\n"
	}

	// Ticket Link
	message += "> *Jira Ticket:* <" + opsLink + "|" + opsName + ">\n"

	messageText := slack.NewTextBlockObject("mrkdwn", message, false, false)

	title := slack.SectionBlock{
		Type:      "section",
		Text:      headerText,
		BlockID:   "",
		Fields:    nil,
		Accessory: nil,
	}

	body := slack.SectionBlock{
		Type:      "section",
		Text:      messageText,
		BlockID:   "",
		Fields:    nil,
		Accessory: nil,
	}

	channel := "C0HTE8DD3" //ops
	//channel := "C028T30QQU8" //andy-tejero-dev
	api := slack.New(slackConfiguration.ClientSecret)

	_, _, err := api.PostMessage(channel, slack.MsgOptionBlocks(title, body))

	if err != nil {
		return err
	}
	return nil
}

func doSlackNotifications(submission ticketRequest, createdIssue *jira.Issue) error {

	err := sendSlackMessage(submission.ticket.Reporter, "Your Jira ticket was successfully created, and assigned to "+submission.ticket.Assignee+", you can find it here -> "+jiraInstance.Config[0].URL+"/browse/"+createdIssue.Key)

	if err != nil {
		botLog(err, "error sending jira issue created confirmation from slack to issue reporter ["+submission.ticket.Assignee+"]")
		return err
	}

	if submission.ticket.Assignee != "N/A" {

		err = sendSlackMessage(submission.ticket.Assignee, "You've been assigned a Jira priority request from "+submission.ticket.Reporter+", you can find it here -->"+jiraInstance.Config[0].URL+"/browse/"+createdIssue.Key)

		if err != nil {
			botLog(err, "error sending jira issue created confirmation from slack to issue assignee ["+submission.ticket.Assignee+"]")
			return err
		}
	}

	err = postMessageToOps(submission, createdIssue.Key, jiraInstance.Config[0].URL+"/browse/"+createdIssue.Key)

	if err != nil {
		botLog(err, "error posting message to #ops")
		return err
	}

	return nil
}
