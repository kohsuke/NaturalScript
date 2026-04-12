package main

import (
	"fmt"
	"os"
	"strings"
)

// Print converts a Script structure back into the file format.
func Print(s Script) (string, error) {
	compressedPrompt, err := Encode([]byte(s.Prompt))
	if err != nil {
		return "", err
	}

	shebang := s.Shebang
	if shebang == "" {
		executable, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("failed to obtain the executable path: %w", err)
		}
		shebang = "#!" + executable
	}

	return strings.Join([]string{shebang, s.Prompt, compressedPrompt, s.GeneratedCode}, Separator), nil
}
