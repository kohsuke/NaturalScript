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

	shebang := s.Shebang
	if shebang == "" {
		shebang = "#!/bin/genscript"
	}

	return fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s\n%s",
		shebang,
		s.Prompt,
		Separator,
		compressedPrompt,
		Separator,
		s.GeneratedCode), nil
}
