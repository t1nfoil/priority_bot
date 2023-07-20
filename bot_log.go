package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
)

func botLog(thing interface{}, message ...string) {

	fmt.Println()
	if len(message) != 0 {
		var max int = 0
		for _, v := range message {
			n, _ := fmt.Printf("[%v] %v\n", time.Now().Format("2006-01-02 15:04:05"), v)
			if n > max {
				max = n
			}
		}
		fmt.Println(strings.Repeat("-", max))
	}

	if botConfiguration.Debug == "true" {
		if thing != nil {
			switch thing.(type) {
			case slack.SlashCommand:
				processedSlackCommand, err := detokenizeSlackLogging(thing)
				if err != nil {
					spew.Dump("Error unmarshalling slack object to detokenize", err)
					spew.Dump(processedSlackCommand)
				}
				spew.Dump(processedSlackCommand)
			default:
				spew.Dump(thing)
			}
		}

	}
}

func detokenizeSlackLogging(thing interface{}) (*slack.SlashCommand, error) {
	data, err := json.Marshal(thing)

	if err != nil {
		return &slack.SlashCommand{}, err
	}

	slackObject := &slack.SlashCommand{}
	json.Unmarshal(data, slackObject)

	slackObject.Token = "************************************** [redacted]"
	return slackObject, nil
}
