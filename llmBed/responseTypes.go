package llmBed

/*
* Json format taken from doc
 */
type GPTRequest struct {
	Model    string        `json:"model"`
	Messages []GPTMessage  `json:"messages"`
	MaxTokens int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`

}


type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}


type GPTResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []struct {
		Message      GPTMessage `json:"message"`
		FinishReason string     `json:"finish_reason"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}