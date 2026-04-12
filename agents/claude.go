package agents

import (
	"os"
	"os/exec"
)

type ClaudeAgent struct{}

func NewClaudeAgent() *ClaudeAgent {
	return &ClaudeAgent{}
}

func (a *ClaudeAgent) Run(prompt string) error {
	cmd := exec.Command("claude", prompt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
