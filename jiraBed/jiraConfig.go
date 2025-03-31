package jiraBed

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func CreateNewIssue(){

}

func JiraClient(username string, password string, url string) *jira.Client{
	jt := jira.BasicAuthTransport{
        Username: username,
        Password: password,
    }

    client, err := jira.NewClient(jt.Client(), url)
    if err != nil {
        fmt.Println(err)
    }

    me, _, err := client.User.GetSelf()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(me)

    return client
}