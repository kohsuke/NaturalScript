package main

type Script struct {
	// the current prompt, which might have been modified by the user since the last generation
	Prompt string
	// the frozen prompt that was used to generate GeneratedCode
	CapturedPrompt string
	// the generated script
	GeneratedCode string
}

const PromptBeginMarker = "==== NATURALSCRIPT:BEGIN ===="
const PromptEndMarker = "==== NATURALSCRIPT:END ===="

func (script Script) ShouldRegenerate() bool {
	if script.GeneratedCode == "" {
		return true
	} else if script.CapturedPrompt != script.Prompt {
		return true
	}
	return false
}
