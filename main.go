package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"genscript/agents"
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

	parts, parseErr := Parse(string(content))

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
		outFile, err := os.CreateTemp("", "genscript-output-*.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temporary output file: %v\n", err)
			os.Exit(1)
		}
		outPath := outFile.Name()
		if err := outFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error preparing temporary output file: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(outPath)

		selectedAgent := os.Getenv("GENSCRIPT_AGENT")
		var a agents.Agent
		switch selectedAgent {
		case "opencode":
			a = agents.NewOpenCodeAgent()
		case "claude":
			a = agents.NewClaudeAgent()
		default:
			a = agents.NewOpenCodeAgent()
		}

		fullPrompt := fmt.Sprintf(
			"Current instruction:\n%s\n\nPrevious generated script (if any):\n%s\n\nGenerate an updated executable script that satisfies the current instruction. Write only the script source code to this exact file path, overwriting it if needed:\n%s\n\nWhen done writing the file, exit.",
			parts.Prompt,
			parts.GeneratedCode,
			outPath,
		)

		fmt.Println("Prompt change detected. Invoking agent...")
		err = a.Run(fullPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Agent error: %v\n", err)
			os.Exit(1)
		}

		newCodeBytes, err := os.ReadFile(outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading generated script from %s: %v\n", outPath, err)
			os.Exit(1)
		}
		newCode := strings.TrimSpace(string(newCodeBytes))
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

		fullFile, err := Print(parts)
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

	if err := Execute(parts, args); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if code := exitErr.ExitCode(); code >= 0 {
				os.Exit(code)
			}
		}
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}
}
