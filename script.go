package main

type Script struct {
	Shebang string
	// the current prompt, which might have been modified by the user since the last generation
	Prompt string
	// the frozen prompt that was used to generate GeneratedCode
	CapturedPrompt string
	// the generated script
	GeneratedCode string
}

const Separator = "-=-=-=-=-=-=-=-="
