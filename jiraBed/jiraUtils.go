/*
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                 │
│ *  NOTE : - For right now gpt is not giving result we take the result we got from analysing the pod its         │
│ separated we told gpt                                                                                           │
│ to                                                                                                              │
│ *        format our response like this :                                                                        │
│ *                     Format your response as <'ISSUE: <brief description>' followed by 'SOLUTION:>             │
│ *                                                                                                               │
│ *                                                                                                               │
│ *  TODO : - We should let the user decide if we let gpt make the ticket or if it takes only what is generated   │
│ from anylisis user should add inputs because its note nough *                                                   │
│ information                                                                                                     │
│                                                                                                           │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
*/
package jiraBed

import (
	"fmt"
	"strings"

	"github.com/G00dBunny/cloud-bunny/llmBed"
	"github.com/andygrunwald/go-jira"
)


type TicketDetail struct {
	Summary string
	Description string
	IssueType string
	Priority string
}

/*
*	DONE
*/
func GenerateTicketFromAnalysis (result llmBed.AnalysisResult) TicketDetail {
	var issue, solution string
	if result.Analysis != "" {
		parts := strings.Split(result.Analysis, "SOLUTION:")
		if len(parts) > 0 {
			issueParts := strings.Split(parts[0], "ISSUE:")
			if len(issueParts) > 1 {
				issue = strings.TrimSpace(issueParts[1])
			}
		}
		if len(parts) > 1 {
			solution = strings.TrimSpace(parts[1])
		}
	}

	return TicketDetail{
		Summary:     fmt.Sprintf("Pod Issue: %s in namespace %s - %s", result.PodName, result.Namespace, issue),
		Description: fmt.Sprintf("Pod: %s\nNamespace: %s\n\nISSUE:\n%s\n\nSOLUTION:\n%s",
			result.PodName, result.Namespace, issue, solution),
		IssueType:   "[System] Incident",
		Priority:    "Medium",
	}
}


/*
* NOTE : provided by https://github.com/andygrunwald/go-jira
*/
func CreateTicket(client *jira.Client, projectKey string, ticketDetail TicketDetail) (*jira.Issue, error) {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     ticketDetail.Summary,
			Description: ticketDetail.Description,
			Project: jira.Project{
				Key: projectKey,
			},
			Type: jira.IssueType{
				Name: ticketDetail.IssueType, 
			},
		},
	}
	
	if ticketDetail.Priority != "" {
		i.Fields.Priority = &jira.Priority{
			Name: ticketDetail.Priority,
		}
	}
	
	issue, resp, err := client.Issue.Create(&i)
	if err != nil {
		return nil, fmt.Errorf("failed to create jira ticket: %v, response: %v", err, resp)
	}
	
	return issue, nil
}