package llmBed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/G00dBunny/cloud-bunny/listutils"
	"k8s.io/client-go/kubernetes"
)

const (
	GPTEndpoint = "https://api.openai.com/v1/chat/completions"
)

type Configuration struct {
	APIKey     string
	Model      string
	MaxTokens  int
	TimeoutSec int
}

type AnalysisResult struct {
	PodName   string
	Namespace string
	Logs      string
	Analysis  string
	Error     error
}

func NewConfig(apiKey, model string, maxTokens, timeoutSec int) Configuration {
	return Configuration{
		APIKey:     apiKey,
		Model:      model,
		MaxTokens:  maxTokens,
		TimeoutSec: timeoutSec,
	}
}




func processLogs(logs string) string {

	if logs == "" {
		return "No logs available"
	}


	if len(logs) > 1000 {
		logs = logs[len(logs)-1000:]
	}

	var relevantLines []string
	lines := strings.Split(logs, "\n")
	
	keyPhrases := []string{
		"error", "Error", "ERROR",
		"warn", "Warn", "WARNING", "Warning",
		"fatal", "Fatal", "FATAL",
		"critical", "Critical", "CRITICAL",
		"exception", "Exception", "EXCEPTION",
		"fail", "Fail", "FAIL",
		"crash", "Crash", "CRASH",
		"OOM", "Out of memory", "out of memory",
		"kill", "Kill", "KILL",
		"permission", "Permission", "PERMISSION",
	}

	for _, line := range lines {
		for _, phrase := range keyPhrases {
			if strings.Contains(line, phrase) {
				relevantLines = append(relevantLines, line)
				break
			}
		}
	}

	if len(relevantLines) == 0 && len(lines) > 0 {
		startIdx := 0
		if len(lines) > 10 {
			startIdx = len(lines) - 10
		}
		relevantLines = lines[startIdx:]
	}

	return strings.Join(relevantLines, "\n")
}

func AnalyzeBadPods(namespaces []string, clientset *kubernetes.Clientset, config Configuration) []AnalysisResult {
	badPods := listutils.GetBadPod(namespaces, clientset)
	results := make([]AnalysisResult, 0, len(badPods))
	
	for _, podName := range badPods {

		podNamespace := listutils.FindPodNamespace(podName, namespaces, clientset)
		if podNamespace == "" {
			results = append(results, AnalysisResult{
				PodName: podName,
				Error:   fmt.Errorf("could not determine namespace for pod %s", podName),
			})
			continue
		}

		logs := listutils.GetPodLog(podNamespace, podName, clientset)
		
		processedLogs := processLogs(logs)
		
		analysis, err := sendToGPT(podName, podNamespace, processedLogs, config)
		
		results = append(results, AnalysisResult{
			PodName:   podName,
			Namespace: podNamespace,
			Logs:      logs,
			Analysis:  analysis,
			Error:     err,
		})
	}
	
	return results
}


/*
*
*	GPT suggestion from 2023 doc 
*
*/
func sendToGPT(podName, namespace, logs string, config Configuration) (string, error) {
    if config.APIKey == "" {
        return "", fmt.Errorf("OpenAI API key not provided")
    }

    prompt := fmt.Sprintf(
        "Analyze these Kubernetes logs from pod '%s' in namespace '%s'. Prioritize in this order: 1) Memory issues 2) Crashes 3) Permission errors. Provide only the critical issue and a simple, practical solution:\n\n%s",
        podName, namespace, logs,
    )

    gptReq := GPTRequest{
        Model: config.Model,
        Messages: []GPTMessage{
            {
                Role:    "system",
                Content: "You are a Kubernetes emergency responder. Be extremely concise. Format your response as 'ISSUE: <brief description>' followed by 'SOLUTION: <practical fix>'. Focus on immediate solutions a human can implement, not commands.",
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        MaxTokens:   config.MaxTokens,
        Temperature: 0.0, 
    }
    
    jsonData, err := json.Marshal(gptReq)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %v", err)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.TimeoutSec)*time.Second)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "POST", GPTEndpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %v", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("API request failed: %v", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %v", err)
    }
    
    var gptResp GPTResponse
    if err := json.Unmarshal(body, &gptResp); err != nil {
        return "", fmt.Errorf("failed to parse response: %v", err)
    }
    
    if gptResp.Error.Message != "" {
        return "", fmt.Errorf("API error: %s", gptResp.Error.Message)
    }
    
    if len(gptResp.Choices) > 0 {
        return gptResp.Choices[0].Message.Content, nil
    }
    
    return "", fmt.Errorf("no response from GPT API")
}