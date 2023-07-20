package main

import (
	"github.com/slack-go/slack"
)

// The function below, is probably the ugliest piece of code you'll see all week. It's gross. It's not elegant, streamlined or
// whatever. It's untidy, unorganized and should be refactored. That being said, it's pathfinding code. Look how we tie the
// action_id's to the block_id's, and you'll see that we can probably make a struct to handle toggling of action id's.. that's
// the next step: turn the code below into a nice library of functions for working with modal blocks (and a small caching backend)

func generateModal(toggleRollStatus OPS_ROLL_STATUS, toggleDeployFlavour OPS_JIRA_DEPLOY_FLAVOUR, toggleProject OPS_JIRA_PROJECT) slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "Priority Ops Request", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit Ops Ticket", false, false)

	// priority
	priorityPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select a ticket priority", false, false)
	priorityOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, priorityPlaceholder, "a_id_priority",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_PRIORITY_P1_DESC)}, Value: string(OPS_PRIORITY_P1)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_PRIORITY_P2_DESC)}, Value: string(OPS_PRIORITY_P2)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_PRIORITY_P3_DESC)}, Value: string(OPS_PRIORITY_P3)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_PRIORITY_P4_DESC)}, Value: string(OPS_PRIORITY_P4)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_PRIORITY_P5_DESC)}, Value: string(OPS_PRIORITY_P5)},
	)
	priorityLabel := slack.NewTextBlockObject(slack.PlainTextType, "Ticket Priority", false, false)
	priorityLevel := slack.NewInputBlock("b_id_priority", priorityLabel, priorityOptions)

	rollStatusText := slack.NewTextBlockObject(slack.MarkdownType, "*Roll Type*", false, false)
	rollStatusPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select Rollback or Rollout", false, false)
	rollStatusOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, rollStatusPlaceholder, "a_id_rollstatus",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_ROLL_INITIAL)}, Value: string(OPS_ROLL_INITIAL)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_ROLL_BACK)}, Value: string(OPS_ROLL_BACK)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_ROLL_OUT)}, Value: string(OPS_ROLL_OUT)},
	)

	rollStatus := slack.NewActionBlock("b_id_rollstatus", rollStatusOptions)
	rollStatusSection := slack.NewSectionBlock(rollStatusText, nil, nil)

	// toggle the roll out or roll back script blocks.
	var rollScriptText *slack.TextBlockObject
	var rollScriptPlaceholder *slack.TextBlockObject
	var rollScriptElement *slack.PlainTextInputBlockElement
	var rollScriptLink *slack.InputBlock

	var rollScriptTestedText *slack.TextBlockObject
	var rollScriptTestedPlaceholder *slack.TextBlockObject
	var rollScriptTestedOptions *slack.SelectBlockElement
	var rollScriptTestedStatus *slack.InputBlock

	if toggleRollStatus == OPS_ROLL_BACK {
		rollScriptText = slack.NewTextBlockObject(slack.PlainTextType, "Roll Back Script", false, false)
		rollScriptPlaceholder = slack.NewTextBlockObject(slack.PlainTextType, "Link the Roll Back script hosted in a Git Repo", false, false)
		rollScriptElement = slack.NewPlainTextInputBlockElement(rollScriptPlaceholder, "a_id_rollbackscript")
		rollScriptLink = slack.NewInputBlock("b_id_rollbackscript", rollScriptText, rollScriptElement)
		rollScriptLink.Optional = true

		rollScriptTestedText = slack.NewTextBlockObject(slack.PlainTextType, "Has the Roll Back script been tested?", false, false)
		rollScriptTestedPlaceholder = slack.NewTextBlockObject(slack.PlainTextType, "Select a Roll Back script test status", false, false)
		rollScriptTestedOptions = slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, rollScriptTestedPlaceholder, "a_id_rollbacktest",
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_YES)}, Value: string(OPS_TESTED_YES)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_NO)}, Value: string(OPS_TESTED_NO)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_UNSURE)}, Value: string(OPS_TESTED_UNSURE)},
		)
		rollScriptTestedStatus = slack.NewInputBlock("b_id_rollbacktest", rollScriptTestedText, rollScriptTestedOptions)
	}
	if toggleRollStatus == OPS_ROLL_OUT {
		rollScriptText = slack.NewTextBlockObject(slack.PlainTextType, "Roll Out Script", false, false)
		rollScriptPlaceholder = slack.NewTextBlockObject(slack.PlainTextType, "Link the Roll Out script hosted in a Git Repo", false, false)
		rollScriptElement := slack.NewPlainTextInputBlockElement(rollScriptPlaceholder, "a_id_rolloutscript")
		rollScriptLink = slack.NewInputBlock("b_id_rolloutscript", rollScriptText, rollScriptElement)
		rollScriptLink.Optional = true

		rollScriptTestedText := slack.NewTextBlockObject(slack.PlainTextType, "Has the Roll Out script been tested?", false, false)
		rollScriptTestedPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select a Roll Out script test status", false, false)
		rollScriptTestedOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, rollScriptTestedPlaceholder, "a_id_rollouttest",
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_YES)}, Value: string(OPS_TESTED_YES)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_NO)}, Value: string(OPS_TESTED_NO)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_TESTED_UNSURE)}, Value: string(OPS_TESTED_UNSURE)},
		)
		rollScriptTestedStatus = slack.NewInputBlock("b_id_rollouttest", rollScriptTestedText, rollScriptTestedOptions)
	}

	// request

	requestText := slack.NewTextBlockObject(slack.PlainTextType, "Request Title", false, false)
	requestPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Type the request title for the Jira ticket", false, false)
	requestElement := slack.NewPlainTextInputBlockElement(requestPlaceholder, "a_id_request")
	requestTitle := slack.NewInputBlock("b_id_request", requestText, requestElement)

	// details

	detailsText := slack.NewTextBlockObject(slack.PlainTextType, "Details", false, false)
	detailsPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the details of the request here", false, false)
	detailsElement := slack.NewPlainTextInputBlockElement(detailsPlaceholder, "a_id_details")
	detailsElement.Multiline = true
	details := slack.NewInputBlock("b_id_details", detailsText, detailsElement)

	// reason

	reasonText := slack.NewTextBlockObject(slack.PlainTextType, "Reason for Change", false, false)
	reasonPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the reason for this request", false, false)
	reasonElement := slack.NewPlainTextInputBlockElement(reasonPlaceholder, "a_id_reason")
	reasonElement.Multiline = true
	reason := slack.NewInputBlock("b_id_reason", reasonText, reasonElement)

	// impact

	impactText := slack.NewTextBlockObject(slack.PlainTextType, "Impact of Change", false, false)
	impactPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the expected impact of this change", false, false)
	impactElement := slack.NewPlainTextInputBlockElement(impactPlaceholder, "a_id_impact")
	impactElement.Multiline = true
	impact := slack.NewInputBlock("b_id_impact", impactText, impactElement)

	// impact level

	impactOutageText := slack.NewTextBlockObject(slack.PlainTextType, "Impact Outage Time", false, false)
	impactOutagePlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the expected impact outage time", false, false)
	impactOutageElement := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, impactOutagePlaceholder, "a_id_impactoutage",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_UNKNOWN)}, Value: string(OPS_IMPACT_UNKNOWN)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_NONE)}, Value: string(OPS_IMPACT_NONE)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_CUSTOMER)}, Value: string(OPS_IMPACT_CUSTOMER)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_1MIN)}, Value: string(OPS_IMPACT_1MIN)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_5MIN)}, Value: string(OPS_IMPACT_5MIN)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_10MIN)}, Value: string(OPS_IMPACT_10MIN)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_30MIN)}, Value: string(OPS_IMPACT_30MIN)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_IMPACT_OVER30MIN)}, Value: string(OPS_IMPACT_OVER30MIN)},
	)

	impactOutage := slack.NewInputBlock("b_id_impactoutage", impactOutageText, impactOutageElement)

	// region

	regionPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select affected Product and Regions (multi-select)", false, false)
	regionText := slack.NewTextBlockObject(slack.PlainTextType, "Select Affected Region(s)", false, false)

	gridRegionObjects := []*slack.OptionGroupBlockObject{
		{
			Label: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Regions", Emoji: false, Verbatim: false},
			Options: []*slack.OptionBlockObject{
				{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "N/A"}, Value: "N/A"},
		},
	}

	var regionElements slack.MultiSelectBlockElement

	regionElements.OptionGroups = gridRegionObjects
	regionElements.Placeholder = regionPlaceholder
	regionElements.Type = slack.MultiOptTypeStatic
	regionElements.ActionID = "a_id_region"

	region := slack.NewInputBlock("b_id_region", regionText, regionElements)

	// due date

	duedateText := slack.NewTextBlockObject(slack.PlainTextType, "Due Date", false, false)
	duedateElement := slack.NewDatePickerBlockElement("a_id_duedate")
	duedate := slack.NewInputBlock("b_id_duedate", duedateText, duedateElement)

	// jira project

	jiraProjectPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select Project Type", false, false)
	jiraProjectOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, jiraProjectPlaceholder, "a_id_jira_project",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_PROJECT_OPS_DESC)}, Value: string(OPS_JIRA_PROJECT_OPS)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_PROJECT_OPSWORK_DESC)}, Value: string(OPS_JIRA_PROJECT_OPSWORK)},
	)
	jiraProject := slack.NewActionBlock("b_id_jira_project", jiraProjectOptions)

	// jira ticket type

	jiraTicketPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select the Jira Ticket Type", false, false)
	var jiraTicketOptions *slack.SelectBlockElement

	if toggleProject == OPS_JIRA_PROJECT_OPS {
		jiraTicketOptions = slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, jiraTicketPlaceholder, "a_id_jira_ticket",
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_TICKET_TYPE_DEPLOY)}, Value: string(OPS_JIRA_TICKET_TYPE_DEPLOY)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_TICKET_TYPE_MAINTENANCE)}, Value: string(OPS_JIRA_TICKET_TYPE_MAINTENANCE)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_TICKET_TYPE_TASK)}, Value: string(OPS_JIRA_TICKET_TYPE_TASK)},
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_TICKET_TYPE_EPIC)}, Value: string(OPS_JIRA_TICKET_TYPE_EPIC)},
		)
	}

	if toggleProject == OPS_JIRA_PROJECT_OPSWORK {
		jiraTicketOptions = slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, jiraTicketPlaceholder, "a_id_jira_ticket",
			&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_TICKET_TYPE_TASK)}, Value: string(OPS_JIRA_TICKET_TYPE_TASK)},
		)
	}

	jiraTicket := slack.NewActionBlock("b_id_jira_ticket", jiraTicketOptions)

	// deployment flavour -- only appears if the action block for Jira Ticket Type shows
	jiraDeploymentFlavourText := slack.NewTextBlockObject(slack.PlainTextType, "Is this first time deployment, or re-deployment?", false, false)
	jiraDeploymentFlavourOptions := slack.NewRadioButtonsBlockElement("a_id_jira_deployment_flavour",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_DEPLOY_FLAVOUR_FIRST)}, Value: string(OPS_JIRA_DEPLOY_FLAVOUR_FIRST)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY)}, Value: string(OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY)},
	)
	jiraDeploymentFlavour := slack.NewInputBlock("b_id_jira_deployment_flavour", jiraDeploymentFlavourText, jiraDeploymentFlavourOptions)

	// jira redeployment link existing ticket

	jiraRedeploymentLinkText := slack.NewTextBlockObject(slack.PlainTextType, "Link to an existing JIRA ticket", false, false)
	jiraRedeploymentLinkPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Enter the Jira link", false, false)
	jiraRedeploymentLinkElement := slack.NewPlainTextInputBlockElement(jiraRedeploymentLinkPlaceholder, "a_id_redeployment_link")
	jiraRedeploymentLink := slack.NewInputBlock("b_id_redeployment_link", jiraRedeploymentLinkText, jiraRedeploymentLinkElement)
	jiraRedeploymentLink.Optional = true

	// assignee
	assigneeText := slack.NewTextBlockObject(slack.PlainTextType, "Select Ticket Assignee", false, false)
	assigneePlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Typically Brendan or Shimon", false, false)
	assigneeOptions := slack.SelectBlockElement{
		Type:        slack.OptTypeUser,
		Placeholder: assigneePlaceholder,
		ActionID:    "a_id_assignee",
		//		InitialUser: "U01PRCB92QK",
	}
	assignee := slack.NewInputBlock("b_id_assignee", assigneeText, assigneeOptions)
	assignee.Optional = true

	// Infrastructure Threat
	infrastructureThreatText := slack.NewTextBlockObject(slack.PlainTextType, "Does this Request affect region Infrastructure?", false, false)
	infrastructureThreatPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Select Yes, No or Unsure", false, false)
	infrastructureThreatOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, infrastructureThreatPlaceholder, "a_id_infrastructure_threat",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_INFRASTRUCTURE_THREAT_YES)}, Value: string(OPS_INFRASTRUCTURE_THREAT_YES)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_INFRASTRUCTURE_THREAT_NO)}, Value: string(OPS_INFRASTRUCTURE_THREAT_NO)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_INFRASTRUCTURE_THREAT_UNSURE)}, Value: string(OPS_INFRASTRUCTURE_THREAT_UNSURE)},
	)
	infrastructureThreat := slack.NewInputBlock("b_id_infrastructure_threat", infrastructureThreatText, infrastructureThreatOptions)

	// CAB requirement
	cabRequiredText := slack.NewTextBlockObject(slack.PlainTextType, "Does this Request require CAB approval?", false, false)
	cabRequiredPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "If unsure, select Yes", false, false)
	cabRequiredOptions := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, cabRequiredPlaceholder, "a_id_cab_required",
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_REQUIRES_CAB_YES)}, Value: string(OPS_REQUIRES_CAB_YES)},
		&slack.OptionBlockObject{Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: string(OPS_REQUIRES_CAB_NO)}, Value: string(OPS_REQUIRES_CAB_NO)},
	)
	cabRequired := slack.NewInputBlock("b_id_cab_required", cabRequiredText, cabRequiredOptions)

	headerBlockJira := slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, "Jira Ticket Information", false, false))

	var blocks slack.Blocks

	blocks.BlockSet = append(blocks.BlockSet, headerBlockJira)
	blocks.BlockSet = append(blocks.BlockSet, jiraProject)

	if toggleProject == OPS_JIRA_PROJECT_OPS || toggleProject == OPS_JIRA_PROJECT_OPSWORK {
		blocks.BlockSet = append(blocks.BlockSet, jiraTicket)
	}

	if toggleDeployFlavour != OPS_JIRA_DEPLOY_FLAVOR_INITIAL && toggleProject == OPS_JIRA_PROJECT_OPS {
		blocks.BlockSet = append(blocks.BlockSet, jiraDeploymentFlavour)
		blocks.BlockSet = append(blocks.BlockSet, jiraRedeploymentLink)
	}

	blocks.BlockSet = append(blocks.BlockSet, duedate)
	blocks.BlockSet = append(blocks.BlockSet, priorityLevel)
	blocks.BlockSet = append(blocks.BlockSet, requestTitle)
	blocks.BlockSet = append(blocks.BlockSet, details)

	if toggleDeployFlavour != OPS_JIRA_DEPLOY_FLAVOR_INITIAL && toggleProject == OPS_JIRA_PROJECT_OPS {
		if toggleRollStatus == OPS_ROLL_INITIAL {
			blocks.BlockSet = append(blocks.BlockSet, rollStatusSection)
			blocks.BlockSet = append(blocks.BlockSet, rollStatus)
		} else {
			blocks.BlockSet = append(blocks.BlockSet, rollStatusSection)
			blocks.BlockSet = append(blocks.BlockSet, rollStatus)
			blocks.BlockSet = append(blocks.BlockSet, rollScriptTestedStatus)
			blocks.BlockSet = append(blocks.BlockSet, rollScriptLink)
		}
	}
	blocks.BlockSet = append(blocks.BlockSet, infrastructureThreat)
	blocks.BlockSet = append(blocks.BlockSet, region)
	blocks.BlockSet = append(blocks.BlockSet, reason)
	blocks.BlockSet = append(blocks.BlockSet, impact)
	blocks.BlockSet = append(blocks.BlockSet, impactOutage)
	blocks.BlockSet = append(blocks.BlockSet, cabRequired)
	blocks.BlockSet = append(blocks.BlockSet, assignee)

	var modalRequest slack.ModalViewRequest

	modalRequest.CallbackID = "PRIORITYBOT_VIEW"
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest

}
