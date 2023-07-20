package main

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	psh "github.com/platformsh/gohelper"
)

var botConfiguration botConfig
var slackConfiguration slackConfig
var jiraConfiguration jiraConfig
var opsRequest ticketQueue
var jiraInstance jiraCloud

func main() {

	botLog(nil, "priority_bot starting")
	spew.Config.Indent = "\t"

	platformInfo, err := psh.NewPlatformInfo()

	if err != nil {
		botLog(nil, "error loading platform.sh app info")
		botLog(nil, err.Error())
		return
	}

	botLog(nil, "listening on "+platformInfo.Project+":"+platformInfo.Port)

	err = botConfiguration.loadBotConfiguration("./config/bot_config.json")

	if err != nil {
		botLog(nil, "error loading bot configuration")
		botLog(nil, err.Error())
	}

	err = slackConfiguration.loadSlackConfiguration("./token/slack_tokens.json")

	if err != nil {
		botLog(nil, "error loading slack token(s)")
		botLog(nil, err.Error())
	}

	err = jiraConfiguration.loadJiraConfiguration("./token/jira_tokens.json")

	if err != nil {
		botLog(nil, "error loading jira token(s)")
		botLog(nil, err.Error())
	}

	jiraInstance.setContext()
	jiraInstance.configJiraAuth()
	jiraInstance.getAuthorization()

	http.HandleFunc("/slash", handleSlackSlashCommand)
	http.HandleFunc("/modal", handleSlackModal)
	http.HandleFunc("/jiracallback", handleJiraCallback)

	botLog(nil, "handlers active, listening...")
	http.ListenAndServe(":"+platformInfo.Port, nil)
}
