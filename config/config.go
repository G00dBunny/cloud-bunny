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

	return jiraBed.JiraConfig{
		Username: jiraUser,
		Token:    jiraToken,
		URL:      jiraURL,
		Project:  jiraProject,
	}, nil
}

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