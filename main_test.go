package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShouldRegenerate(t *testing.T) {
	tests := []struct {
		name   string
		script Script
		want   bool
	}{
		{
			name: "regenerates when generated code is empty",
			script: Script{
				Prompt:         "task",
				CapturedPrompt: "task",
				GeneratedCode:  "",
			},
			want: true,
		},
		{
			name: "regenerates when prompt changed",
			script: Script{
				Prompt:         "new task",
				CapturedPrompt: "old task",
				GeneratedCode:  "echo hi",
			},
			want: true,
		},
		{
			name: "does not regenerate when code exists and prompt unchanged",
			script: Script{
				Prompt:         "task",
				CapturedPrompt: "task",
				GeneratedCode:  "echo hi",
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.script.ShouldRegenerate()
			if got != tc.want {
				t.Fatalf("ShouldRegenerate() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	input := []byte("line one\nline two\nline three")

	encoded, err := Encode(input)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if string(decoded) != string(input) {
		t.Fatalf("Decode(Encode(x)) mismatch: got %q, want %q", string(decoded), string(input))
	}
}

func TestDecodeInvalidInput(t *testing.T) {
	if _, err := Decode("this-is-not-valid-base64"); err == nil {
		t.Fatalf("Decode() expected error for invalid base64 input")
	}
}

func TestParseWithShebangCapturedPromptAndGeneratedCode(t *testing.T) {
	prompt := "generate a script"
	generated := "echo first" + Separator + "echo second"

	content, _ := Print(Script{
		Shebang:        "#!/usr/bin/env naturalscript",
		Prompt:         prompt,
		CapturedPrompt: prompt,
		GeneratedCode:  generated,
	})

	s := Parse(content)

	if s.Shebang != "#!/usr/bin/env naturalscript" {
		t.Fatalf("Shebang = %q, want %q", s.Shebang, "#!/usr/bin/env naturalscript")
	}
	if s.Prompt != prompt {
		t.Fatalf("Prompt = %q, want %q", s.Prompt, prompt)
	}
	if s.CapturedPrompt != prompt {
		t.Fatalf("CapturedPrompt = %q, want %q", s.CapturedPrompt, prompt)
	}
	if s.GeneratedCode != generated {
		t.Fatalf("GeneratedCode = %q, want %q", s.GeneratedCode, generated)
	}
}

func TestFormatArguments(t *testing.T) {
	args := []string{"simple", "with space", "quote\"and\\slash"}
	got := formatArguments(args)
	want := `["simple", "with space", "quote\"and\\slash"]`

	if got != want {
		t.Fatalf("formatArguments() = %q, want %q", got, want)
	}
}

func TestPromptForRevision(t *testing.T) {
	s := Script{
		Prompt:         "new task",
		CapturedPrompt: "old task",
		GeneratedCode:  "echo old",
	}
	msg := prompt(s, "/tmp/out.sh", nil)

	if !strings.Contains(msg, s.CapturedPrompt) {
		t.Fatalf("prompt() missing captured prompt")
	}
	if !strings.Contains(msg, s.GeneratedCode) {
		t.Fatalf("prompt() missing generated code")
	}
	if !strings.Contains(msg, s.Prompt) {
		t.Fatalf("prompt() missing new prompt")
	}
}

func TestMakeTmpFileCreatesAFile(t *testing.T) {
	path, err := makeTmpFile()
	if err != nil {
		t.Fatalf("makeTmpFile() error = %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat temp file: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("expected temp file, got directory")
	}
}

func TestAtomicWriteReplacesContentsAndSetsExecutableBit(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "script.sh")

	if err := os.WriteFile(target, []byte("old"), 0644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if err := atomicWrite(target, "new"); err != nil {
		t.Fatalf("atomicWrite() error = %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(data) != "new" {
		t.Fatalf("file content = %q, want %q", string(data), "#!/bin/sh\\necho new")
	}
}
