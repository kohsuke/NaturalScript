package main

import (
	"fmt"
	"strings"
)

const metadataInstruction = "Managed by NaturalScript. Edit the prompt below and run `naturalscript path/to/this/file` to regenerate. DO NOT MODIFY THE REST"

// Print converts a Script structure back into the file format.
func Print(s Script) (string, error) {
	if s.GeneratedCode == "" {
		return "", fmt.Errorf("generated code is empty")
	}

	shebang, body := splitFirstLine(s.GeneratedCode)
	if !strings.HasPrefix(shebang, "#!") {
		return "", fmt.Errorf("expecting #!... but found '%s'", shebang)
	}

	compressedPrompt, err := Encode([]byte(s.Prompt))
	if err != nil {
		return "", err
	}

	envelope := EnvelopeForShebang(shebang)
	metadata := strings.Join([]string{
		metadataInstruction,
		PromptBeginMarker,
		s.Prompt,
		PromptEndMarker,
		compressedPrompt,
		"", // end of base64 marker
	}, "\n")

	return shebang + "\n" + envelope.render(metadata) + "\n" + body, nil
}

type Envelope struct {
	BeginEnvelope string
	PerLinePrefix string
	EndEnvelope   string
}

func EnvelopeForShebang(shebang string) Envelope {
	lower := strings.ToLower(shebang)

	if isShellShebang(lower) {
		return Envelope{
			BeginEnvelope: ": <<'COMMENTBLOCK_FOR_NATURALSCRIPT'",
			EndEnvelope:   "COMMENTBLOCK_FOR_NATURALSCRIPT",
		}
	}

	if strings.Contains(lower, "python") {
		return Envelope{
			BeginEnvelope: "'''",
			EndEnvelope:   "'''",
		}
	}

	// fallback
	return Envelope{PerLinePrefix: "# "}
}

func isShellShebang(lower string) bool {
	shells := []string{"bash", "zsh", "ksh", "dash", "ash", "sh"}
	for _, s := range shells {
		if strings.Contains(lower, "bin/"+s) || strings.Contains(lower, "env "+s) {
			return true
		}
	}
	return false
}

func (e Envelope) render(metadata string) string {
	if e.PerLinePrefix != "" {
		lines := strings.Split(metadata, "\n")
		for i := range lines {
			lines[i] = e.PerLinePrefix + lines[i]
		}
		return strings.Join(lines, "\n")
	}

	return e.BeginEnvelope + "\n" + metadata + "\n" + e.EndEnvelope
}

func splitFirstLine(s string) (string, string) {
	if s == "" {
		return "", ""
	}
	idx := strings.IndexByte(s, '\n')
	if idx == -1 {
		return s, ""
	}
	return s[:idx], s[idx+1:]
}
