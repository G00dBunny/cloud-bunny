package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/G00dBunny/cloud-bunny/jiraBed"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}


/*
* FIXME : this is a temporary and test implementation.... should use types and not this ugly function... 
*/
// func GetJiraConfig() (string, string, string, string) {
// 	username := os.Getenv("JIRA_USER")
// 	token := os.Getenv("JIRA_TOKEN")
// 	url := os.Getenv("JIRA_URL")
// 	project := os.Getenv("JIRA_PROJECT")
	
// 	// createTickets, err := strconv.ParseBool(os.Getenv("CREATE_JIRA_TICKETS"))
// 	// if err != nil {
// 	// 	createTickets = false 
// 	// }
	
// 	return username, token, url, project
// }

/*
*	DONE : use this instead
*/
func GetJiraConfig() (jiraBed.JiraConfig, error) {
	
	jiraUser := os.Getenv("JIRA_USER")
	if jiraUser == "" {
		return jiraBed.JiraConfig{}, errors.New("JIRA_USER environment variable is not set")
	}

	jiraToken := os.Getenv("JIRA_TOKEN")
	if jiraToken == "" {
		return jiraBed.JiraConfig{}, errors.New("JIRA_TOKEN environment variable is not set")
	}


	jiraURL := os.Getenv("JIRA_URL")
	if jiraURL == "" {
		return jiraBed.JiraConfig{}, errors.New("JIRA_URL environment variable is not set")
	}

	jiraProject := os.Getenv("JIRA_PROJECT")
	if jiraProject == "" {
		return jiraBed.JiraConfig{}, errors.New("JIRA_PROJECT environment variable is not set")
	}

	

	// isUsedStr := os.Getenv("CREATE_JIRA_TICKETS")
	// if isUsedStr == "" {
	// 	return &jiraBed.JiraConfig{}, errors.New("Set usage of jira in environment variable ")
	// }

	// isUsed, err := strconv.ParseBool(isUsedStr)
	// if err != nil {
	// 	return &jiraBed.JiraConfig{}, errors.New("CREATE_JIRA_TICKETS must be a valid boolean value")
	// }

	return jiraBed.JiraConfig{
		Username: jiraUser,
		Token:    jiraToken,
		URL:      jiraURL,
		Project:  jiraProject,
	}, nil
}



/*
*	DONE
*/
func GetOpenAIConfig() (string, string, int, int) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	
	model := os.Getenv("GPT_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	
	maxTokensStr := os.Getenv("GPT_MAX_TOKENS")
	maxTokens := 1000
	if maxTokensStr != "" {
		if val, err := strconv.Atoi(maxTokensStr); err == nil {
			maxTokens = val
		}
	}
	
	timeoutStr := os.Getenv("GPT_TIMEOUT_SEC")
	timeout := 30
	if timeoutStr != "" {
		if val, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = val
		}
	}
	
	return apiKey, model, maxTokens, timeout
}