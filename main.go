package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"genscript/agents"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: genscript <script-path> [script-args...]")
		os.Exit(1)
	}

	scriptPath := os.Args[1]
	args := os.Args[2:]

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading script: %v\n", err)
		os.Exit(1)
	}

	script := Parse(string(content))

	if script.ShouldRegenerate() {
		outFile, err := os.CreateTemp("", "genscript-output-*.txt")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating temporary output file: %v\n", err)
			os.Exit(1)
		}
		outPath := outFile.Name()
		if err := outFile.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error preparing temporary output file: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(outPath)

		a := selectAgent()

		fullPrompt := prompt(script, outPath, args)

		err = a.Run(fullPrompt)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Agent error: %v\n", err)
			os.Exit(1)
		}

		newCodeBytes, err := os.ReadFile(outPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading generated script from %s: %v\n", outPath, err)
			os.Exit(1)
		}
		newCode := string(newCodeBytes)
		if newCode == "" {
			fmt.Println("Agent did not produce a script. Exiting.")
			os.Exit(1)
		}

		script.GeneratedCode = newCode
		script.CapturedPrompt = script.Prompt

		if script.Shebang == "" {
			script.Shebang = "#!" + os.Args[0]
		}

		content, err := Print(script)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error serializing script: %v\n", err)
			os.Exit(1)
		}

		err = atomicWrite(scriptPath, content)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error serializing script: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := Execute(script, args); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			_, _ = fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
			os.Exit(1)
		}
	}
}

func atomicWrite(scriptPath string, contents string) error {
	// Write to a temporary file in the same directory as scriptPath
	scriptDir := "."
	if dir := strings.TrimSuffix(scriptPath, "/"+filepath.Base(scriptPath)); dir != "" {
		scriptDir = dir
	}
	tmpFile, err := os.CreateTemp(scriptDir, "genscript-tmp-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	defer tmpFile.Close()

	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return err
	}
	err = tmpFile.Chmod(0755)
	if err != nil {
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	err = os.Rename(tmpPath, scriptPath)
	if err != nil {
		return err
	}

	return nil
}

func prompt(script Script, outPath string, args []string) string {
	var prompt string

	if script.GeneratedCode == "" {
		prompt = fmt.Sprintf(`
I'd like to turn the following repeatable task into a script: 
====
%s
====

`, script.Prompt)
	} else {
		prompt = fmt.Sprintf(`
I wanted to turn the following repeatable task into a script: 
====
%s
====

You earlier gave me the following script for this task:
====
%s
====

Now, my task definition changed into the following:
====
%s
====

I'd like you to produce a revised script that reflects this change. 
`, script.CapturedPrompt, script.GeneratedCode, script.Prompt)
	}

	prompt += fmt.Sprintf(`
In order to produce the correct script, first I'd like you to be the interpreter.
Ask me any clarifying questions, and execute the necessary commands directly.

When we are done, please use that knowledge to write out the script to %s, so that the next time this same task
can be performed without you.

For this session, the "arguments" I'm invoking this script with are: %s
`, outPath, formatArguments(args))

	return prompt
}

func formatArguments(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = "'" + arg + "'"
	}
	arguments := "[" + strings.Join(quoted, ", ") + "]"
	return arguments
}

func selectAgent() agents.Agent {
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
	return a
}
