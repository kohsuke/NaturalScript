package agents

import (
	"os"
	"os/exec"
)

type OpenCodeAgent struct{}

func NewOpenCodeAgent() *OpenCodeAgent {
	return &OpenCodeAgent{}
}

func (a *OpenCodeAgent) Run(prompt string) error {
	cmd := exec.Command("opencode", prompt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
