package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
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