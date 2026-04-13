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
				GeneratedCode:  "#!/bin/sh\necho hi",
			},
			want: true,
		},
		{
			name: "does not regenerate when code exists and prompt unchanged",
			script: Script{
				Prompt:         "task",
				CapturedPrompt: "task",
				GeneratedCode:  "#!/bin/sh\necho hi",
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

func TestPrintAndParsePythonEnvelopeRoundTrip(t *testing.T) {
	in := Script{
		Prompt:        "Say hello in Italian",
		GeneratedCode: "#!/usr/bin/env python3\nprint('Ciao, mondo!')",
	}

	printed, err := Print(in)
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if !strings.Contains(printed, "'''") {
		t.Fatalf("expected python triple-quote envelope")
	}
	if !strings.Contains(printed, PromptBeginMarker) {
		t.Fatalf("expected metadata begin marker")
	}

	out, _ := Parse(printed)
	if out.Prompt != in.Prompt {
		t.Fatalf("Prompt = %q, want %q", out.Prompt, in.Prompt)
	}
	if out.CapturedPrompt != in.Prompt {
		t.Fatalf("CapturedPrompt = %q, want %q", out.CapturedPrompt, in.Prompt)
	}
	if out.GeneratedCode != in.GeneratedCode {
		t.Fatalf("GeneratedCode = %q, want %q", out.GeneratedCode, in.GeneratedCode)
	}
}

func TestPrintAndParseShellHeredocEnvelopeRoundTrip(t *testing.T) {
	in := Script{
		Prompt:        "echo from shell",
		GeneratedCode: "#!/usr/bin/env bash\necho hello",
	}

	printed, err := Print(in)
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}
	if !strings.Contains(printed, ": <<'COMMENTBLOCK_FOR_NATURALSCRIPT'") {
		t.Fatalf("expected shell heredoc envelope")
	}

	out, _ := Parse(printed)
	if out.Prompt != in.Prompt {
		t.Fatalf("Prompt = %q, want %q", out.Prompt, in.Prompt)
	}
	if out.CapturedPrompt != in.Prompt {
		t.Fatalf("CapturedPrompt = %q, want %q", out.CapturedPrompt, in.Prompt)
	}
	if out.GeneratedCode != in.GeneratedCode {
		t.Fatalf("GeneratedCode = %q, want %q", out.GeneratedCode, in.GeneratedCode)
	}
}

func TestParseLineCommentFallbackEnvelope(t *testing.T) {
	compressed, err := Encode([]byte("task2"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	content := strings.Join([]string{
		"#!/usr/bin/env ruby",
		"# " + metadataInstruction,
		"# " + PromptBeginMarker,
		"# task",
		"# " + PromptEndMarker,
		"# " + compressed,
		"# ", // end of base64 block
		"# ", // end of multiline comment
		"puts 'hello'",
	}, "\n")

	s, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if s.Prompt != "task" {
		t.Fatalf("Prompt = %q, want %q", s.Prompt, "task")
	}
	if s.CapturedPrompt != "task2" {
		t.Fatalf("CapturedPrompt = %q, want %q", s.CapturedPrompt, "task2")
	}
	if s.GeneratedCode != "#!/usr/bin/env ruby\nputs 'hello'" {
		t.Fatalf("GeneratedCode = %q", s.GeneratedCode)
	}
}

func TestParseBash(t *testing.T) {
	compressed, err := Encode([]byte("task2"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	content := strings.Join([]string{
		"#!/usr/bin/bash",
		": <<EOF",
		"warning warning",
		PromptBeginMarker,
		"task",
		PromptEndMarker,
		compressed,
		"",    // end of base64 block
		"EOF", // end of multiline comment
		"echo hello",
	}, "\n")

	s, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if s.Prompt != "task" {
		t.Fatalf("Prompt = %q, want %q", s.Prompt, "task")
	}
	if s.CapturedPrompt != "task2" {
		t.Fatalf("CapturedPrompt = %q, want %q", s.CapturedPrompt, "task2")
	}
	if s.GeneratedCode != "#!/usr/bin/bash\necho hello" {
		t.Fatalf("GeneratedCode = %q", s.GeneratedCode)
	}
}

func TestPrintWithoutShebangFails(t *testing.T) {
	if _, err := Print(Script{Prompt: "task", GeneratedCode: "echo ok"}); err == nil {
		t.Fatalf("Print() expected error when generated code lacks shebang")
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
		GeneratedCode:  "#!/bin/sh\necho old",
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
	path, err := makeTmpFile("./foo")
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

	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
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
		t.Fatalf("file content = %q, want %q", string(data), "new")
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat target: %v", err)
	}
	if info.Mode()&0o100 == 0 {
		t.Fatalf("expected owner executable bit to be set, mode = %v", info.Mode())
	}
}
