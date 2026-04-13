package main

import (
	"os"
	"path/filepath"
)

func makeTmpFile(scriptPath string) (string, error) {
	outFile, err := os.CreateTemp(filepath.Dir(scriptPath), "naturalscript-output-*.txt")
	if err != nil {
		return "", err
	}
	outPath := outFile.Name()
	if err := outFile.Close(); err != nil {
		return "", err
	}
	return outPath, nil
}

func atomicWrite(scriptPath string, contents string) error {
	tmpFile, err := os.CreateTemp(filepath.Dir(scriptPath), "naturalscript-tmp-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return err
	}
	err = tmpFile.Chmod(0755)
	if err != nil {
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, scriptPath)
}
