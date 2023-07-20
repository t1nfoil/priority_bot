## priority_bot

Sends your `#ops` priority requests to Jira. Is also used to submit CAB requests for the CAB committee.

---



### Where tokens are stored

You will need to make a tokens/ folder off the root of this project, and create your slack and jira token files (check the comments in slack_config.go and jira_config.go for token file syntax) or set associated environment variales in lieu of the JSON token files.

You can also define and use platform sh specific env variables such as below if you don't like tokens.

| env variable name | sensitive / build / runtime |
| ------ | ------ |
| BOT_JIRA_CLIENT_ID | build, runtime |
| BOT_JIRA_CLIENT_SECRET | sensitive / build / runtime |
| BOT_SLACK_APP_ID | build / runtime |
| BOT_SLACK_CLIENT_ID | build / runtime |
| BOT_SLACK_CLIENT_SECRET |  sensitive / build / runtime |
| BOT_SLACK_SIGNING_SECRET |  sensitive / build / runtime |


### Permissions needed for upstream service(s) operations (Slack, Jira)

Ensure the correct permissions are enabled on the jira account the bot app will be hosted from, in particular these scopes

#### Jira

| permission level | permission scope |
| ------ | ------ |
| manage | jira-webhook |
| write | jira-work |
| read | jira-work |
| write | jira-user |
| read | jira-user |


#### Slack

For the slack permissions, these are the required oauth scopes

| permission level | permission scope | purpose |
| ------ | ------ | ------ |
| read | app_mentions | View messages that directly mention @priority_bot in conversations that the app is in |
| join | channels | Join public channels in a workspace |
| read | channels | View basic information about public channels in a workspace |
| write | channels | Send messages as @priority_bot |
| write | chat | |
| write.customize | chat | Send messages as @priority_bot with a customized username and avatar |
| | commands | Add shortcuts and/or slash commands that people can use |
| | incoming-webhook | Post messages to specific channels in Slack |
| read | links | View URLs in messages |
| read | users.profile | View profile details about people in a workspace |
| read | users | View people in a workspace |
| read.email | users | View email addresses of people in a workspace |
| execute | workflow.steps | Add steps that people can use in Workflow Builder |
    

### TODO 

- lets move this sucker to websockets so we don't have to implement a public endpoint
- add ability to page ops on call from bot when request is super urgent
- write a bot-code generator for future bots

### To request a change or a feature

DM @andy edit this readme and add the feature request below this line.

---


