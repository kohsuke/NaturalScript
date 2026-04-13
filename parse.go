package main

import (
	"fmt"
	"strings"
)

// Parse splits the script file into its constituent parts and decodes the captured prompt.
func Parse(content string) (Script, error) {
	if strings.Contains(content, PromptBeginMarker) {
		return parseManagedScript(content)
	} else {
		// initial hand-written prompt
		return Script{Prompt: content}, nil
	}
}

func parseManagedScript(content string) (Script, error) {
	lines := strings.Split(content, "\n")

	if !strings.HasPrefix(lines[0], "#!") {
		return Script{}, fmt.Errorf("expecting #!... but found '%s'", lines[0])
	}
	shebang := lines[0]

	beginIdx := -1
	prefix := ""
	for i := 1; i < len(lines); i++ {
		pos := strings.Index(lines[i], PromptBeginMarker)
		if pos >= 0 {
			beginIdx = i
			prefix = lines[i][0:pos]
			break
		}
	}
	if beginIdx == -1 {
		panic("did not find metadata begin marker")
	}

	endIdx := -1
	for i := beginIdx + 1; i < len(lines); i++ {
		pos := strings.Index(lines[i], PromptEndMarker)
		if pos >= 0 {
			if lines[i][0:pos] != prefix {
				return Script{}, fmt.Errorf("end marker is incorrectly indented at line %d", i)
			}
			endIdx = i
			break
		}
	}
	if endIdx == -1 {
		return Script{}, fmt.Errorf("end marker '%s' not found", PromptEndMarker)
	}

	promptLines := make([]string, 0, endIdx-beginIdx-1)
	for i := beginIdx + 1; i < endIdx; i++ {
		promptLines = append(promptLines, strings.TrimPrefix(lines[i], prefix))
	}

	// parse base64 section until we see the empty line
	base64EndIdx := -1
	base64Lines := []string{}
	for i := endIdx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], prefix) {
			l := strings.TrimPrefix(lines[i], prefix)
			if l == "" {
				base64EndIdx = i
				break
			}
			base64Lines = append(base64Lines, l)
		} else {
			return Script{}, fmt.Errorf("expecting base64 encoded section but missingat line %d", i)
		}
	}

	codeStart := base64EndIdx + 2 // give one line to wrap up the multiline comment

	code := ""
	if codeStart < len(lines) {
		code = strings.Join(lines[codeStart:], "\n")
	}

	prompt := strings.Join(promptLines, "\n")

	capturedPrompt := ""
	compressed := strings.TrimSpace(strings.Join(base64Lines, "\n"))
	if compressed != "" {
		captured, err := Decode(compressed)
		if err == nil {
			capturedPrompt = string(captured)
		}
	}

	return Script{
		Prompt:         prompt,
		CapturedPrompt: capturedPrompt,
		GeneratedCode:  shebang + "\n" + code,
	}, nil
}
