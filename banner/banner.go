package banner

import (
	"fmt"
	"os"
	"strings"
)

const (
	charHeight    = 8
	firstChar     = 32
	lastChar      = 126
	expectedChars = lastChar - firstChar + 1
)

func Load(name string) ([][]string, error) {
	filename := fmt.Sprintf("banner/%s.txt", name)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open banner file '%s.txt': %w", name, err)
	}

	return parse(data, name)
}

func parse(data []byte, name string) ([][]string, error) {
	content := strings.ReplaceAll(string(data), "\r\n", "\n")

	lines := strings.Split(content, "\n")

	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}

	blockSize := charHeight + 1

	chars := make([][]string, 0, expectedChars)

	for i := 0; i+charHeight <= len(lines); i += blockSize {
		block := make([]string, charHeight)
		for row := 0; row < charHeight; row++ {
			block[row] = lines[i+row]
		}
		chars = append(chars, block)
	}

	if len(chars) < expectedChars {
		return nil, fmt.Errorf(
			"banner '%s' is malformed: expected %d characters, got %d",
			name, expectedChars, len(chars),
		)
	}

	return chars[:expectedChars], nil
}
