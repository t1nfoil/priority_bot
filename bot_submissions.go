package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

type OPS_ASSIGNEE string
type OPS_JIRA_PROJECT string
type OPS_JIRA_TICKET_TYPE string
type OPS_JIRA_DEPLOY_FLAVOUR string
type OPS_PRIORITY string
type OPS_PRODUCT_TYPE string
type OPS_ROLL_STATUS string
type OPS_REGION string
type OPS_TESTED string
type OPS_IMPACT_OUTAGE string
type OPS_REQUIRES_CAB string
type OPS_INFRASTRUCTURE_THREAT string

const (
	OPS_INFRASTRUCTURE_THREAT_YES    OPS_INFRASTRUCTURE_THREAT = "Yes"
	OPS_INFRASTRUCTURE_THREAT_NO     OPS_INFRASTRUCTURE_THREAT = "No"
	OPS_INFRASTRUCTURE_THREAT_UNSURE OPS_INFRASTRUCTURE_THREAT = "Unsure"

	OPS_REQUIRES_CAB_YES OPS_REQUIRES_CAB = "Yes"
	OPS_REQUIRES_CAB_NO  OPS_REQUIRES_CAB = "No"

	OPS_ASSIGNEE_BRENDAN OPS_ASSIGNEE = "@Brendan"
	OPS_ASSIGNEE_SHIMON  OPS_ASSIGNEE = "@Shimon"

	OPS_JIRA_PROJECT_OPS          OPS_JIRA_PROJECT = "11000"
	OPS_JIRA_PROJECT_OPS_DESC     OPS_JIRA_PROJECT = "Ops"
	OPS_JIRA_PROJECT_OPSWORK      OPS_JIRA_PROJECT = "12100"
	OPS_JIRA_PROJECT_OPSWORK_DESC OPS_JIRA_PROJECT = "OPSWORK"

	OPS_JIRA_TICKET_TYPE_DEPLOY      OPS_JIRA_TICKET_TYPE = "Deploy"
	OPS_JIRA_TICKET_TYPE_MAINTENANCE OPS_JIRA_TICKET_TYPE = "Maintenance"
	OPS_JIRA_TICKET_TYPE_TASK        OPS_JIRA_TICKET_TYPE = "Task"
	OPS_JIRA_TICKET_TYPE_EPIC        OPS_JIRA_TICKET_TYPE = "Epic"

	OPS_JIRA_DEPLOY_FLAVOR_INITIAL   OPS_JIRA_DEPLOY_FLAVOUR = ""
	OPS_JIRA_DEPLOY_FLAVOR_SHOW      OPS_JIRA_DEPLOY_FLAVOUR = "Show"
	OPS_JIRA_DEPLOY_FLAVOUR_FIRST    OPS_JIRA_DEPLOY_FLAVOUR = "First Deployment"
	OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY OPS_JIRA_DEPLOY_FLAVOUR = "Re-Deployment"

	OPS_PRIORITY_P1 OPS_PRIORITY = "1"
	OPS_PRIORITY_P2 OPS_PRIORITY = "2"
	OPS_PRIORITY_P3 OPS_PRIORITY = "3"
	OPS_PRIORITY_P4 OPS_PRIORITY = "4"
	OPS_PRIORITY_P5 OPS_PRIORITY = "5"

	OPS_PRIORITY_P1_DESC OPS_PRIORITY = "(P1) Highest - Urgent request and must be dealt with immediately"
	OPS_PRIORITY_P2_DESC OPS_PRIORITY = "(P2) High -  Request that must be dealt with today"
	OPS_PRIORITY_P3_DESC OPS_PRIORITY = "(P3) Medium - Request needing attention fast"
	OPS_PRIORITY_P4_DESC OPS_PRIORITY = "(P4) Low - Request needing attention"
	OPS_PRIORITY_P5_DESC OPS_PRIORITY = "(P5) Lowest - Request to be done when possible"

	OPS_ROLL_INITIAL OPS_ROLL_STATUS = "N/A"
	OPS_ROLL_BACK    OPS_ROLL_STATUS = "Roll Back"
	OPS_ROLL_OUT     OPS_ROLL_STATUS = "Roll Out"

	OPS_TESTED_NO     OPS_TESTED = "No"
	OPS_TESTED_UNSURE OPS_TESTED = "Unsure"
	OPS_TESTED_YES    OPS_TESTED = "Yes"

	OPS_IMPACT_UNKNOWN   OPS_IMPACT_OUTAGE = "Unknown"
	OPS_IMPACT_NONE      OPS_IMPACT_OUTAGE = "None"
	OPS_IMPACT_CUSTOMER  OPS_IMPACT_OUTAGE = "Slowed Performance (No Downtime)"
	OPS_IMPACT_1MIN      OPS_IMPACT_OUTAGE = "Up to 1 min of downtime"
	OPS_IMPACT_5MIN      OPS_IMPACT_OUTAGE = "Up to 5 mins of downtime"
	OPS_IMPACT_10MIN     OPS_IMPACT_OUTAGE = "Up to 10 mins of downtime"
	OPS_IMPACT_30MIN     OPS_IMPACT_OUTAGE = "Up to 30 mins of downtime"
	OPS_IMPACT_OVER30MIN OPS_IMPACT_OUTAGE = "Over 30 mins of downtime"
)

type opsTicket struct {
	MetaTicketTime       time.Time                 `json:"meta_ticket_time"`
	MetaTicketUUID       uuid.UUID                 `json:"meta_ticket_uuid"`
	Reporter             string                    `json:"reporter"`
	Assignee             string                    `json:"assignee"`
	CabRequired          OPS_REQUIRES_CAB          `json:"cab_required"`
	Details              string                    `json:"details"`
	DueDate              string                    `json:"duedate"`
	DeployFlavour        OPS_JIRA_DEPLOY_FLAVOUR   `json:"deploy_flavour"`
	RedeployJiraLink     string                    `json:"redeploy_jira_link"`
	JiraProject          OPS_JIRA_PROJECT          `json:"jira_project"`
	JiraTicketType       OPS_JIRA_TICKET_TYPE      `json:"jira_ticket_type"`
	ChangeReason         string                    `json:"change_reason"`
	Impact               string                    `json:"impact"`
	ImpactOutage         OPS_IMPACT_OUTAGE         `json:"impact_outage"`
	InfrastructureThreat OPS_INFRASTRUCTURE_THREAT `json:"infrastructure_threat"`
	PriorityLevel        OPS_PRIORITY              `json:"priority_level"`
	ProductType          OPS_PRODUCT_TYPE          `json:"product_type"`
	Regions              []string                  `json:"regions"`
	RequestTitle         string                    `json:"request_title"`
	RollBackScriptLink   string                    `json:"roll_back_script_link"`
	RollBackScriptTested OPS_TESTED                `json:"roll_back_script_tested"`
	RollOutScriptLink    string                    `json:"roll_out_script_link"`
	RollOutScriptTested  OPS_TESTED                `json:"roll_out_script_tested"`
	RollStatus           OPS_ROLL_STATUS           `json:"roll_out_status"`
}

func (o *opsTicket) SetCabRequired(cabRequired OPS_REQUIRES_CAB) {
	o.CabRequired = cabRequired
}

func (o *opsTicket) SetInfrastructureThreat(infrastructureThreat OPS_INFRASTRUCTURE_THREAT) {
	o.InfrastructureThreat = infrastructureThreat
}

func (o *opsTicket) SetRedeployJiraLink(redeployJiraLink string) {
	o.RedeployJiraLink = redeployJiraLink
}

func (o *opsTicket) SetAssignee(assignee string) {
	o.Assignee = assignee
}

func (o *opsTicket) SetChangeReason(changeReason string) {
	o.ChangeReason = changeReason
}

func (o *opsTicket) SetImpact(impact string) {
	o.Impact = impact
}

func (o *opsTicket) SetImpactOutage(impactOutage OPS_IMPACT_OUTAGE) {
	o.ImpactOutage = impactOutage
}

func (o *opsTicket) SetReporter(requester string) {
	o.Reporter = requester
}

func (o *opsTicket) SetDetails(details string) {
	o.Details = details
}

func (o *opsTicket) SetDeployFlavour(flavour OPS_JIRA_DEPLOY_FLAVOUR) {
	o.DeployFlavour = flavour
}

func (o *opsTicket) SetJiraProject(jiraProject OPS_JIRA_PROJECT) {
	o.JiraProject = jiraProject
}

func (o *opsTicket) SetJiraTicketType(jiraTicketType OPS_JIRA_TICKET_TYPE) {
	o.JiraTicketType = jiraTicketType
}

func (o *opsTicket) SetPriorityLevel(priorityLevel OPS_PRIORITY) {
	o.PriorityLevel = priorityLevel
}

func (o *opsTicket) SetProductType(productType OPS_PRODUCT_TYPE) {
	o.ProductType = productType
}

func (o *opsTicket) SetRegions(regions []string) {
	o.Regions = regions
}

func (o *opsTicket) SetRequestTitle(requestTitle string) {
	o.RequestTitle = requestTitle
}

func (o *opsTicket) SetRollBackScriptLink(rollBackScriptLink string) {
	o.RollBackScriptLink = rollBackScriptLink
}

func (o *opsTicket) SetRollBackScriptTested(rollBackScriptTested OPS_TESTED) {
	o.RollBackScriptTested = rollBackScriptTested
}
func (o *opsTicket) SetRollOutScriptLink(rollOutScriptLink string) {
	o.RollOutScriptLink = rollOutScriptLink
}

func (o *opsTicket) SetRollOutScriptTested(rollOutScriptTested OPS_TESTED) {
	o.RollOutScriptTested = rollOutScriptTested
}

func (o *opsTicket) SetRollStatus(rollStatus OPS_ROLL_STATUS) {
	o.RollStatus = rollStatus
}

func (o *opsTicket) SetDudeDate(dueDate string) {
	o.DueDate = dueDate
}

func (o *opsTicket) NewTicket() opsTicket {
	return opsTicket{}
}

func (o *opsTicket) PopulateTicket(slackCallBack *slack.InteractionCallback) error {

	for k, v := range slackCallBack.View.State.Values {
		switch k {

		case "b_id_jira_ticket":
			value := v["a_id_jira_ticket"].SelectedOption.Value

			if value == string(OPS_JIRA_TICKET_TYPE_DEPLOY) {
				o.SetJiraTicketType(OPS_JIRA_TICKET_TYPE_DEPLOY)
				break
			}

			if value == string(OPS_JIRA_TICKET_TYPE_MAINTENANCE) {
				o.SetJiraTicketType(OPS_JIRA_TICKET_TYPE_MAINTENANCE)
				break
			}

			if value == string(OPS_JIRA_TICKET_TYPE_TASK) {
				o.SetJiraTicketType(OPS_JIRA_TICKET_TYPE_TASK)
				break
			}

			if value == string(OPS_JIRA_TICKET_TYPE_EPIC) {
				o.SetJiraTicketType(OPS_JIRA_TICKET_TYPE_EPIC)
				break
			}

			return errors.New("error setting jira ticket type - deploy, maintenance or task not specified")

		case "b_id_rollstatus":
			value := v["a_id_rollstatus"].SelectedOption.Value

			if value == string(OPS_ROLL_BACK) {
				o.SetRollStatus(OPS_ROLL_BACK)
				break
			}

			if value == string(OPS_ROLL_OUT) {
				o.SetRollStatus(OPS_ROLL_OUT)
				break
			}

			if value == string(OPS_ROLL_INITIAL) {
				o.SetRollStatus(OPS_ROLL_INITIAL)
				break
			}

			return errors.New("error setting rollstatus - rollout or rollback not specified")

		case "b_id_rollbacktest":

			value := v["a_id_rollbacktest"].SelectedOption.Value

			if value == string(OPS_TESTED_NO) {
				o.SetRollBackScriptTested(OPS_TESTED_NO)
				break
			}

			if value == string(OPS_TESTED_UNSURE) {
				o.SetRollBackScriptTested(OPS_TESTED_UNSURE)
				break
			}

			if value == string(OPS_TESTED_YES) {
				o.SetRollBackScriptTested(OPS_TESTED_YES)
				break
			}

			return errors.New("error setting rollbacktest - no, unsure or yes not specified")

		case "b_id_rollbackscript":

			value := v["a_id_rollbackscript"].Value
			o.SetRollBackScriptLink(value)

		case "b_id_rollouttest":

			value := v["a_id_rollouttest"].SelectedOption.Value

			if value == string(OPS_TESTED_NO) {
				o.SetRollOutScriptTested(OPS_TESTED_NO)
				break
			}

			if value == string(OPS_TESTED_UNSURE) {
				o.SetRollOutScriptTested(OPS_TESTED_UNSURE)
				break
			}

			if value == string(OPS_TESTED_YES) {
				o.SetRollOutScriptTested(OPS_TESTED_YES)
				break
			}

			return errors.New("error setting rollouttest - no, unsure or yes not specified")

		case "b_id_rolloutscript":

			value := v["a_id_rolloutscript"].Value
			o.SetRollOutScriptLink(value)

		case "b_id_priority":

			value := v["a_id_priority"].SelectedOption.Value

			if value == string(OPS_PRIORITY_P1) {
				o.SetPriorityLevel(OPS_PRIORITY_P1)
				break
			}

			if value == string(OPS_PRIORITY_P2) {
				o.SetPriorityLevel(OPS_PRIORITY_P2)
				break
			}

			if value == string(OPS_PRIORITY_P3) {
				o.SetPriorityLevel(OPS_PRIORITY_P3)
				break
			}

			if value == string(OPS_PRIORITY_P4) {
				o.SetPriorityLevel(OPS_PRIORITY_P4)
				break
			}

			if value == string(OPS_PRIORITY_P5) {
				o.SetPriorityLevel(OPS_PRIORITY_P5)
				break
			}

			return errors.New("error setting priority level - P1,P2,P3,P4 or P5 not specified")

		case "b_id_request":

			value := v["a_id_request"].Value

			if value == "" {
				return errors.New("error setting request title - request title not specified")
			}

			o.SetRequestTitle(value)

		case "b_id_details":

			value := v["a_id_details"].Value

			if value == "" {
				return errors.New("error setting request details - details not specified")
			}

			o.SetDetails(value)

		case "b_id_jira_project":

			value := v["a_id_jira_project"].SelectedOption.Value

			if value == string(OPS_JIRA_PROJECT_OPS) {
				o.SetJiraProject(OPS_JIRA_PROJECT_OPS)
				break
			}

			if value == string(OPS_JIRA_PROJECT_OPSWORK) {
				o.SetJiraProject(OPS_JIRA_PROJECT_OPSWORK)
				break
			}

			return errors.New("error setting jira project - type not specified (Ops?)")

		case "b_id_assignee":

			value := v["a_id_assignee"].SelectedUser

			if value == "" {
				botLog(nil, "assignee - not specified")
				o.SetAssignee("N/A")
				break
			}

			botLog(value, "Ticket Assignee Raw Value")

			api := slack.New(slackConfiguration.ClientSecret)
			user, err := api.GetUserInfo(value)

			if err != nil {
				botLog(user, "error user profile email empty in PopulateTicket()")
				return errors.New("error retrieving slack user information for assignee")
			}

			o.SetAssignee(user.Profile.Email)
			botLog(user.Profile.Email, "Ticket Assignee Email")

		case "b_id_jira_deployment_flavour":

			value := v["a_id_jira_deployment_flavour"].SelectedOption.Value

			if value == string(OPS_JIRA_DEPLOY_FLAVOUR_FIRST) {
				o.SetDeployFlavour(OPS_JIRA_DEPLOY_FLAVOUR_FIRST)
				break
			}

			if value == string(OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY) {
				o.SetDeployFlavour(OPS_JIRA_DEPLOY_FLAVOUR_REDEPLOY)
				break
			}

			return errors.New("error setting deployment type - first/redeploy not specified")

		case "b_id_duedate":

			value := v["a_id_duedate"].SelectedDate
			o.SetDudeDate(value)

		case "b_id_region":

			value := v["a_id_region"].SelectedOptions
			var regions []string

			for _, v := range value {
				regions = append(regions, v.Value)
			}

			o.SetRegions(regions)

		case "b_id_impact":

			value := v["a_id_impact"].Value
			o.SetImpact(value)

		case "b_id_impactoutage":

			value := v["a_id_impactoutage"].SelectedOption.Value

			if value == string(OPS_IMPACT_UNKNOWN) {
				o.SetImpactOutage(OPS_IMPACT_UNKNOWN)
				break
			}

			if value == string(OPS_IMPACT_NONE) {
				o.SetImpactOutage(OPS_IMPACT_NONE)
				break
			}

			if value == string(OPS_IMPACT_NONE) {
				o.SetImpactOutage(OPS_IMPACT_NONE)
				break
			}

			if value == string(OPS_IMPACT_CUSTOMER) {
				o.SetImpactOutage(OPS_IMPACT_CUSTOMER)
				break
			}
			if value == string(OPS_IMPACT_1MIN) {
				o.SetImpactOutage(OPS_IMPACT_1MIN)
				break
			}

			if value == string(OPS_IMPACT_5MIN) {
				o.SetImpactOutage(OPS_IMPACT_5MIN)
				break
			}

			if value == string(OPS_IMPACT_10MIN) {
				o.SetImpactOutage(OPS_IMPACT_10MIN)
				break
			}

			if value == string(OPS_IMPACT_30MIN) {
				o.SetImpactOutage(OPS_IMPACT_30MIN)
				break
			}

			if value == string(OPS_IMPACT_OVER30MIN) {
				o.SetImpactOutage(OPS_IMPACT_OVER30MIN)
				break
			}

			return errors.New("error setting impact outage level - no impact level specified")

		case "b_id_reason":

			value := v["a_id_reason"].Value
			o.SetChangeReason(value)

		case "b_id_redeployment_link":

			value := v["a_id_redeployment_link"].Value
			o.SetRedeployJiraLink(value)

		case "b_id_cab_required":

			value := v["a_id_cab_required"].SelectedOption.Value

			if value == string(OPS_REQUIRES_CAB_YES) {
				o.SetCabRequired(OPS_REQUIRES_CAB_YES)
				break
			}

			if value == string(OPS_REQUIRES_CAB_NO) {
				o.SetCabRequired(OPS_REQUIRES_CAB_NO)
				break
			}

			return errors.New("error setting CAB requirement - no CAB requirement specified")

		case "b_id_infrastructure_threat":

			value := v["a_id_infrastructure_threat"].SelectedOption.Value

			if value == string(OPS_INFRASTRUCTURE_THREAT_YES) {
				o.SetInfrastructureThreat(OPS_INFRASTRUCTURE_THREAT_YES)
				break
			}

			if value == string(OPS_INFRASTRUCTURE_THREAT_NO) {
				o.SetInfrastructureThreat(OPS_INFRASTRUCTURE_THREAT_NO)
				break
			}

			if value == string(OPS_INFRASTRUCTURE_THREAT_UNSURE) {
				o.SetInfrastructureThreat(OPS_INFRASTRUCTURE_THREAT_UNSURE)
				break
			}

			return errors.New("error setting infrastructure threat - no infrastructure threat was specified")

		default:
		}

	}

	return nil
}
