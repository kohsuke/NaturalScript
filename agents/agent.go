package agents

// Agent defines the interface for interacting with different LLM agents.
type Agent interface {
	// Run starts an interactive session with the fully prepared prompt.
	Run(prompt string) error
}
