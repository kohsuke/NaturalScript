package main

import (
	"strings"
)

// Parse splits the script file into its constituent parts and decodes the captured prompt.
func Parse(content string) Script {
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
		return s
	}

	s.Prompt = strings.TrimSpace(parts[0])

	if len(parts) >= 2 {
		compressed := strings.TrimSpace(parts[1])
		if compressed != "" {
			captured, err := Decode(compressed)
			if err == nil {
				s.CapturedPrompt = strings.TrimSpace(string(captured))
			}
			// if the decode fails, say due to corruption, then  leave it empty to trigger the regeneration
		}
	}

	if len(parts) >= 3 {
		s.GeneratedCode = strings.Join(parts[2:], Separator)
	}

	return s
}
