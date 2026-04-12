package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute runs the generated script with the provided arguments.
func Execute(s Script, args []string) error {
	if s.GeneratedCode == "" {
		return fmt.Errorf("generated script is empty")
	}

	tmpFile, err := os.CreateTemp("", "naturalscript-*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(s.GeneratedCode); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return err
	}

	cmd := exec.Command(tmpFile.Name(), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
