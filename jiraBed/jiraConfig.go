package jiraBed

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)


type JiraConfig struct {
	Username string
	Token string
	URL string
	Project string
}



/*
*   NOTE : https://github.com/andygrunwald/go-jira
*   DONE
*/
func JiraClient(config JiraConfig) (*jira.Client, error){
	jt := jira.BasicAuthTransport{
        Username: config.Username,
        Password: config.Token,
    }

    client, err := jira.NewClient(jt.Client(), config.URL)
    if err != nil {
    return nil, fmt.Errorf("jira client fail: %v", err)
    }

    _, _, err = client.User.GetSelf()
    if err != nil {
		return nil, fmt.Errorf("jira auth fail: %v", err)
    }

    return client, nil
}