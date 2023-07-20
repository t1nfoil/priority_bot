package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func getJiraListOfIssueTypes(jiraClient *jira.Client) ([]jira.IssueType, error) {

	apiEndpoint := "/rest/api/2/issuetype"
	req, err := jiraClient.NewRequestWithContext(jiraInstance.ctx, "GET", apiEndpoint, nil)

	if err != nil {
		return nil, err
	}

	issueTypes := []jira.IssueType{}
	_, err = jiraClient.Do(req, &issueTypes)

	if err != nil {
		return nil, err
	}

	return issueTypes, nil
}

func getJiraListOfIssueFields(jiraClient *jira.Client) ([]jira.Field, error) {

	fields, _, err := jiraClient.Field.GetList()
	if err != nil {
		return nil, fmt.Errorf("unable to get the list of fields: %w", err)
	}

	fieldMapping := make(map[string]string)

	for _, f := range fields {
		if f.Custom {
			fieldMapping[f.Name] = f.ID
		}
	}

	return fields, err
}

func getJiraListOfProjects(jiraClient *jira.Client) (*jira.ProjectList, error) {
	projects, _, err := jiraClient.Project.GetListWithContext(jiraInstance.ctx)

	if err != nil {
		return nil, err
	}

	return projects, err
}

func getJiraListOfPriorities(jiraClient *jira.Client) ([]jira.Priority, error) {
	priority, _, err := jiraClient.Priority.GetListWithContext(jiraInstance.ctx)

	if err != nil {
		return nil, err
	}

	return priority, err
}

func getJiraUserByEmail(jiraClient *jira.Client, emailAddress string) ([]jira.User, error) {
	user, _, err := jiraClient.User.FindWithContext(jiraInstance.ctx, emailAddress)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getAllJiraUsers(jiraClient *jira.Client) ([]jira.User, error) {
	apiEndpoint := "/rest/api/2/users/search?maxResults=1000"
	req, err := jiraClient.NewRequestWithContext(jiraInstance.ctx, "GET", apiEndpoint, nil)

	if err != nil {
		return nil, err
	}

	users := []jira.User{}
	_, err = jiraClient.Do(req, &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}
