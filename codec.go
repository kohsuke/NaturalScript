package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

const base64LineWidth = 78

// Encode compresses data using gzip and then encodes it in base64.
func Encode(data []byte) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	if len(encoded) <= base64LineWidth {
		return encoded, nil
	}

	var out strings.Builder
	for i := 0; i < len(encoded); i += base64LineWidth {
		end := i + base64LineWidth
		if end > len(encoded) {
			end = len(encoded)
		}
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(encoded[i:end])
	}

	return out.String(), nil
}

// Decode decodes a base64 string and then decompresses it using gzip.
func Decode(encoded string) ([]byte, error) {
	normalized := strings.Join(strings.Fields(encoded), "")
	data, err := base64.StdEncoding.DecodeString(normalized)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return io.ReadAll(gz)
}
