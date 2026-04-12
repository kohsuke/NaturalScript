package internal

import (
	"strings"
	"genscript/internal/codec"
)


// Parse splits the script file into its constituent parts and decodes the captured prompt.
func Parse(content string) ScriptParts {
	lines := strings.Split(content, "\n")
	startIdx := 0
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!") {
		startIdx = 1
	}
	
	remainingContent := strings.Join(lines[startIdx:], "\n")
	parts := strings.Split(remainingContent, Separator)
	
	if len(parts) == 1 {
		return ScriptParts{
			Prompt: strings.TrimSpace(parts[0]),
		}
	}
	
	if len(parts) == 2 {
		prompt := strings.TrimSpace(parts[0])
		compressed := strings.TrimSpace(parts[1])
		captured, _ := codec.Decode(compressed)
		return ScriptParts{
			Prompt:         prompt,
			CapturedPrompt: string(captured),
		}
	}
	
	prompt := strings.TrimSpace(parts[0])
	compressed := strings.TrimSpace(parts[1])
	captured, _ := codec.Decode(compressed)
	
	return ScriptParts{
		Prompt:         prompt,
		CapturedPrompt: string(captured),
		GeneratedCode:  strings.TrimSpace(strings.Join(parts[2:], Separator)),
	}
}
