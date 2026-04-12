package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"genscript/internal"
	"genscript/internal/agent"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: #!/bin/genscript")
		os.Exit(1)
	}

	scriptPath := os.Args[0]
	args := os.Args[1:]

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading script: %v\n", err)
		os.Exit(1)
	}

	parts, parseErr := internal.Parse(string(content))

	shouldRecompile := false
	if parts.GeneratedCode == "" {
		shouldRecompile = true
	} else if parseErr != nil {
		shouldRecompile = true
	} else if strings.TrimSpace(parts.CapturedPrompt) != parts.Prompt {
		shouldRecompile = true
	}

	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v; triggering regeneration\n", parseErr)
	}

	if shouldRecompile {
		selectedAgent := os.Getenv("GENSCRIPT_AGENT")
		var a agent.Agent
		switch selectedAgent {
		case "opencode":
			a = agent.NewOpenCodeAgent()
		case "claude":
			a = agent.NewClaudeAgent()
		default:
			a = agent.NewOpenCodeAgent()
		}

		fmt.Println("Prompt change detected. Invoking agent...")
		newCode, err := a.Run(parts.Prompt, parts.GeneratedCode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Agent error: %v\n", err)
			os.Exit(1)
		}

		if newCode == "" {
			fmt.Println("Agent did not produce a script. Exiting.")
			os.Exit(1)
		}

		parts.GeneratedCode = newCode

		// Infer shebang from current executable if missing
		if parts.Shebang == "" {
			exe, err := os.Executable()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
				os.Exit(1)
			}
			parts.Shebang = "#!" + exe
		}

		fullFile, err := internal.Print(parts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error serializing script: %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(scriptPath, []byte(fullFile), 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating script file: %v\n", err)
			os.Exit(1)
		}
	}

	if err := internal.Execute(parts, args); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if code := exitErr.ExitCode(); code >= 0 {
				os.Exit(code)
			}
		}
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}
}
