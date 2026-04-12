package main

import (
	"fmt"
	"os"
)

// Print converts a Script structure back into the file format.
func Print(s Script) (string, error) {
	compressedPrompt, err := Encode([]byte(s.Prompt))
	if err != nil {
		return "", err
	}

	if s.Shebang == "" {
		executable, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("obtain executable path: %w", err)
		}
		s.Shebang = "#!" + executable
	}

	return fmt.Sprintf("%s\n%s%s%s%s%s", s.Shebang, s.Prompt, Separator, compressedPrompt, Separator, s.GeneratedCode), nil
}
