package main

import (
	"fmt"
)

// Print converts a Script structure back into the file format.
func Print(s Script) (string, error) {
	compressedPrompt, err := Encode([]byte(s.Prompt))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", s.Shebang, s.Prompt, Separator, compressedPrompt, Separator, s.GeneratedCode), nil
}
