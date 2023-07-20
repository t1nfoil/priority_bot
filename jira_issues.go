package main

import (
	"errors"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/trivago/tgo/tcontainer"
)

func populateIssueReporterAndAssignee(jiraClient *jira.Client, submission ticketRequest) (jira.User, jira.User, error) {
	reporterAccountID, err := getJiraUserByEmail(jiraClient, submission.ticket.Reporter)

	if err != nil {
		botLog(submission, "error getting the jira account ID for the submission reporting user")
		return jira.User{}, jira.User{}, err
	}

	if submission.ticket.Assignee == "N/A" {
		return reporterAccountID[0], jira.User{}, nil
	}

	assigneeAccountID, err := getJiraUserByEmail(jiraClient, submission.ticket.Assignee)

	if err != nil {
		botLog(submission, "error getting jira account ID for assigned user in submission ticket")
		return jira.User{}, jira.User{}, err
	}

	return reporterAccountID[0], assigneeAccountID[0], err
}

func populateIssueCommentBlock(submission ticketRequest) string {
	var commentBlock string

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_DEPLOY {
		commentBlock += "*Issue Type:* Deployment\n"
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_TASK {
		commentBlock += "*Issue Type:* Task\n"
	}
	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_MAINTENANCE {
		commentBlock += "*Issue Type:* Maintenance\n"
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_EPIC {
		commentBlock += "*Issue Type:* Epic\n"
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_EPIC {
		commentBlock += "*Issue Type:* Story\n"
	}

	if submission.ticket.DeployFlavour == OPS_JIRA_DEPLOY_FLAVOUR_FIRST {
		commentBlock += "*Deployment Type:* First Deployment\n"
	}

	if submission.ticket.DeployFlavour == OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY {
		commentBlock += "*Deployment Type:* Redeployment\n"
		commentBlock += "*Linked Issues:* " + submission.ticket.RedeployJiraLink + "\n"
	}

	commentBlock += "*Request affects Region Infrastructure:* " + string(submission.ticket.InfrastructureThreat) + "\n\n\n"

	commentBlock += "*Affected Region(s):*\n"

	for _, v := range submission.ticket.Regions {
		commentBlock += "* " + v + "\n"
	}

	commentBlock += "\n\n"

	if submission.ticket.RollStatus == OPS_ROLL_BACK {
		commentBlock += "*Roll type:* Rollback\n"
		commentBlock += "*Rollback script tested:* " + string(submission.ticket.RollBackScriptTested) + "\n"

		if submission.ticket.RollBackScriptLink != "" {
			commentBlock += "*Rollback script link:* " + string(submission.ticket.RollBackScriptLink) + "\n"
		}
	}

	if submission.ticket.RollStatus == OPS_ROLL_OUT {
		commentBlock += "*Roll type:* Rollout\n"
		commentBlock += "*Rollout script tested:* " + string(submission.ticket.RollOutScriptTested) + "\n"

		if submission.ticket.RollOutScriptLink != "" {
			commentBlock += "*Rollout script link:* " + string(submission.ticket.RollOutScriptLink) + "\n"
		}
	}

	commentBlock += "\n\n"

	commentBlock += "*Reason for Change:* " + submission.ticket.ChangeReason + "\n"
	commentBlock += "*Impact:* " + submission.ticket.Impact + "\n"
	commentBlock += "*Expected Impact Time:* " + string(submission.ticket.ImpactOutage) + "\n"

	return commentBlock
}

func populateIssueType(submission ticketRequest) (jira.IssueType, error) {

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_EPIC {
		return jira.IssueType{ID: "10000"}, nil
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_DEPLOY {
		return jira.IssueType{ID: "11421"}, nil
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_MAINTENANCE {
		return jira.IssueType{ID: "11520"}, nil
	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_TASK && submission.toggleProject == OPS_JIRA_PROJECT_OPS {
		return jira.IssueType{ID: "3"}, nil

	}

	if submission.ticket.JiraTicketType == OPS_JIRA_TICKET_TYPE_TASK && submission.toggleProject == OPS_JIRA_PROJECT_OPSWORK {
		return jira.IssueType{ID: "3"}, nil
	}

	return jira.IssueType{}, errors.New("undefined issue type in submission request [" + string(submission.ticket.JiraTicketType) + "]")
}

func populateJiraFields(submission ticketRequest, issueType jira.IssueType, commentBlock string, assignee, reporter jira.User, dueDate time.Time) jira.IssueFields {
	var fields jira.IssueFields

	if submission.toggleProject == OPS_JIRA_PROJECT_OPSWORK {
		fields = jira.IssueFields{
			Priority:    &jira.Priority{ID: string(submission.ticket.PriorityLevel)},
			Project:     jira.Project{ID: string(submission.ticket.JiraProject)},
			Assignee:    &jira.User{AccountID: assignee.AccountID},
			Reporter:    &jira.User{AccountID: reporter.AccountID},
			Duedate:     jira.Date(dueDate),
			Type:        issueType,
			Summary:     submission.ticket.RequestTitle,
			Description: commentBlock + "\n\n" + submission.ticket.Details,
		}
	} else {
		fields = jira.IssueFields{
			Priority:    &jira.Priority{ID: string(submission.ticket.PriorityLevel)},
			Project:     jira.Project{ID: string(submission.ticket.JiraProject)},
			Assignee:    &jira.User{AccountID: assignee.AccountID},
			Reporter:    &jira.User{AccountID: reporter.AccountID},
			Duedate:     jira.Date(dueDate),
			Type:        issueType,
			Summary:     submission.ticket.RequestTitle,
			Description: commentBlock + "\n\n" + submission.ticket.Details,
		}
	}

	if submission.ticket.Assignee == "N/A" {
		fields.Assignee = nil
	}

	return fields
}

func convertSlackToJiraTime(submission ticketRequest) (time.Time, error) {
	var dueDate time.Time
	var err error

	if len(submission.ticket.DueDate) == 10 {
		reconstructedDate := submission.ticket.DueDate[5:7] + "-" + submission.ticket.DueDate[8:10] + "-" + submission.ticket.DueDate[0:4]
		dueDate, err = time.Parse("01-02-2006", reconstructedDate)

		if err != nil {
			botLog(submission, "error setting jira issue due date from slack date ["+submission.ticket.DueDate+"]")
			return time.Time{}, err
		}

		return dueDate, nil
	}

	return time.Time{}, errors.New("slack date format is incorrect length (> or < 10)")
}

func createIssue(jiraClient *jira.Client, submission ticketRequest) (*jira.Issue, error) {

	commentBlock := populateIssueCommentBlock(submission)
	issueType, err := populateIssueType(submission)

	if err != nil {
		botLog(submission, "error populating issue type")
		return nil, err
	}

	reporter, assignee, err := populateIssueReporterAndAssignee(jiraClient, submission)

	if err != nil {
		botLog(submission, "error populating issue assignee or reporter")
		return nil, err
	}

	dueDate, err := convertSlackToJiraTime(submission)

	if err != nil {
		botLog(submission, "error converting slack date format to jira date format")
		return nil, err
	}

	var issue jira.Issue

	fields := populateJiraFields(submission, issueType, commentBlock, assignee, reporter, dueDate)
	issue.Fields = &fields

	// Custom Fields are applied below

	customFields := tcontainer.NewMarshalMap()

	// customfield_13347 is CAB REQUIRED
	// 13992 - Checkmark Yes
	// 13993 - Checkmark No

	// customfield_13343 is CAB STATUS

	if submission.ticket.CabRequired == OPS_REQUIRES_CAB_YES {
		customFields["customfield_13354"] = map[string]string{"value": "Yes - Required"}
		customFields["customfield_13343"] = map[string]string{"value": "CAB REVIEW"}
	} else {
		customFields["customfield_13354"] = map[string]string{"value": "No - Not Required"}
	}

	issue.Fields.Unknowns = customFields

	createdIssue, _, err := jiraClient.Issue.CreateWithContext(jiraInstance.ctx, &issue)

	if err != nil {
		botLog(issue, "error creating issue: "+err.Error())
		return nil, err
	}
	return createdIssue, nil
}
