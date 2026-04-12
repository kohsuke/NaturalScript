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
	// there appears to be no equivalent of claude -p ...
	// programmatic prompt specification requires "run PROMPT", then we can use -c to launch an interactive session
	cmd := exec.Command("opencode", "run", prompt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
	}

	cmd = exec.Command("opencode", "-c")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
