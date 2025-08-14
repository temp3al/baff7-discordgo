package baff_ollama

type OllamaAPICall struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	StreamResponse bool   `json:"stream"`
}

type OllamaAPIResponse struct {
	Response string `json:"response"`
	// used when StreamResponse is true
	Done bool `json:"done"`
}
