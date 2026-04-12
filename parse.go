package main

import (
	"fmt"
	"strings"
)

// Parse splits the script file into its constituent parts and decodes the captured prompt.
func Parse(content string) (Script, error) {
	lines := strings.Split(content, "\n")
	var s Script
	startIdx := 0
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!") {
		s.Shebang = lines[0]
		startIdx = 1
	}

	remainingContent := strings.Join(lines[startIdx:], "\n")
	parts := strings.Split(remainingContent, Separator)

	if len(parts) == 0 {
		return s, nil
	}

	s.Prompt = strings.TrimSpace(parts[0])

	if len(parts) >= 2 {
		compressed := strings.TrimSpace(parts[1])
		if compressed != "" {
			captured, err := Decode(compressed)
			if err != nil {
				return s, fmt.Errorf("failed to decode captured prompt: %w", err)
			}
			s.CapturedPrompt = string(captured)
		}
	}

	if len(parts) >= 3 {
		s.GeneratedCode = strings.TrimSpace(strings.Join(parts[2:], Separator))
	}

	return s, nil
}
