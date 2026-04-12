package main

import (
	"fmt"
	"os"
	"strings"

	"genscript/internal"
	"genscript/internal/agent"
	"genscript/internal/codec"
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

	parts := internal.Parse(string(content))
	
	shouldRecompile := false
	if parts.GeneratedCode == "" {
		shouldRecompile = true
	} else if strings.TrimSpace(parts.CapturedPrompt) != parts.Prompt {
		shouldRecompile = true
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

		compressedPrompt, err := codec.Encode([]byte(parts.Prompt))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding prompt: %v\n", err)
			os.Exit(1)
		}

		fullFile := fmt.Sprintf("#!/bin/genscript\n%s\n\n%s\n%s\n%s\n%s", 
			parts.Prompt, 
			internal.Separator, 
			compressedPrompt, 
			internal.Separator, 
			newCode)

		err = os.WriteFile(scriptPath, []byte(fullFile), 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating script file: %v\n", err)
			os.Exit(1)
		}
		
		parts.GeneratedCode = newCode
	}

	if err := internal.Execute(parts.GeneratedCode, args); err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}
}
