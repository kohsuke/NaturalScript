package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"naturalscript/agents"

	"github.com/manifoldco/promptui"
)

func main() {
	exitCode, err := run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func run() (int, error) {
	if len(os.Args) < 2 {
		return 1, fmt.Errorf("usage: naturalscript <script-path> [script-args...]")
	}

	scriptPath := os.Args[1]
	args := os.Args[2:]

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return 1, fmt.Errorf("failed to read script: %w", err)
	}

	script, err := Parse(string(content))
	if err != nil {
		return 1, fmt.Errorf("failed to parse script: %w", err)
	}

	if script.ShouldRegenerate() {
		fmt.Println("Triggering the agent to generate the script...")

		outPath, err := makeTmpFile(scriptPath)
		if err != nil {
			return 1, fmt.Errorf("can't create temporary output file: %w", err)
		}
		os.Remove(outPath)       // if a file exists, agent needs to read it first
		defer os.Remove(outPath) // delete the file the agent will create, too

		a, err := selectAgent()
		if err != nil {
			return 1, err
		}
		fullPrompt := prompt(script, outPath, args)
		if err := a.Run(fullPrompt); err != nil {
			return 1, fmt.Errorf("agent error: %w", err)
		}

		newCodeBytes, err := os.ReadFile(outPath)
		if err != nil {
			return 1, fmt.Errorf("read generated script from %s: %w", outPath, err)
		}
		newCode := string(newCodeBytes)
		if newCode == "" {
			return 1, errors.New("agent did not produce a script")
		}

		script.GeneratedCode = newCode
		script.CapturedPrompt = script.Prompt

		serializedScript, err := Print(script)
		if err != nil {
			return 1, fmt.Errorf("failed to write script: %w", err)
		}
		if err := atomicWrite(scriptPath, serializedScript); err != nil {
			return 1, fmt.Errorf("failed to write script: %w", err)
		}
		fmt.Printf("Updated %s\n", scriptPath)
		return 0, nil
	} else {
		// if the script path is just a bare name like 'foo', exec by default looks for PATH,
		// but what we really want is ./foo
		if !strings.Contains(scriptPath, "/") {
			scriptPath = "./" + scriptPath
		}

		cmd := exec.Command(scriptPath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exitErr.ExitCode(), nil
			}
			return 1, fmt.Errorf("failed to run the script: %w", err)
		}
		return 0, nil
	}
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
Assume reasonable defaults, but ask me clarifying questions, and execute the necessary commands directly.

When we are done, please use that knowledge to write out the script to %s, so that
the next time this same task can be performed without you. Unless I change my mind,
assume a shell script.

Important: include a shebang line at the top of the generated script.

Then ask the user to exit the session.

For this session, the "arguments" I'm invoking this script with are: %s
`, outPath, formatArguments(args))

	return prompt
}

func formatArguments(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = strconv.Quote(arg)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func selectAgent() (agents.Agent, error) {
	prompt := promptui.Select{
		Label: "Select agent to generate the script with",
		Items: []string{"OpenCode", "Claude Code"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("no agent selected")
	}
	switch result {
	case "OpenCode":
		return agents.NewOpenCodeAgent(), nil
	case "Claude Code":
		return agents.NewClaudeAgent(), nil
	default:
		panic(result)
	}
}
