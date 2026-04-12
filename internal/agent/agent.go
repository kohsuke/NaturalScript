package agent

// Agent defines the interface for interacting with different LLM agents.
type Agent interface {
	// Run starts the interactive session with the agent.
	// It takes the current prompt and the previous script (if any) as context.
	// It returns the generated script once the user is satisfied.
	Run(prompt string, oldScript string) (string, error)
}
