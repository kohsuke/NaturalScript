package internal

import (
	"os"
	"os/exec"
)

// Execute runs the generated script with the provided arguments.
func Execute(scriptCode string, args []string) error {
	tmpFile, err := os.CreateTemp("", "genscript-*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(scriptCode); err != nil {
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
