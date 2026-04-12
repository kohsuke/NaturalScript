package agent

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ClaudeAgent struct{}

func NewClaudeAgent() *ClaudeAgent {
	return &ClaudeAgent{}
}

func (a *ClaudeAgent) Run(prompt string, oldScript string) (string, error) {
	fullPrompt := fmt.Sprintf("User Request: %s\n\nPrevious Script (if any):\n%s\n\nInstruction: Implement the request. Once the user is satisfied, output the final script enclosed in <<<GENSCRIPT_START>>> and <<<GENSCRIPT_END>>> markers.", prompt, oldScript)

	cmd := exec.Command("claude", fullPrompt)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var scriptBuilder strings.Builder
	inScriptBlock := false
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Echo to user

		if strings.Contains(line, "<<<GENSCRIPT_START>>>") {
			inScriptBlock = true
			continue
		}
		if strings.Contains(line, "<<<GENSCRIPT_END>>>") {
			inScriptBlock = false
			continue
		}

		if inScriptBlock {
			scriptBuilder.WriteString(line + "\n")
		}
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return scriptBuilder.String(), nil
}
