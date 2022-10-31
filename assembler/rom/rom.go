package rom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// ROM contains the loaded ROM dump.
type ROM struct {
	Filename string
	Bytes    []byte
}

func Load(filename string) (ROM, error) {
	rawInput, err := ioutil.ReadFile(filename)
	if err != nil {
		return ROM{}, err
	}
	// fmt.Println(string(rawInput))
	parsedBytes, err := parse(string(rawInput))
	if err != nil {
		return ROM{}, err
	}
	return ROM{
		Filename: filename,
		Bytes:    parsedBytes,
	}, nil
}

func parse(input string) ([]byte, error) {
	var result bytes.Buffer
	lines := strings.Split(input, "\n")
	lines = lines[1:]
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed: %s", line)
		}
		hexes := strings.Split(parts[1], " ")
		if len(hexes) != 16 {
			return nil, fmt.Errorf("must contain 16 hexes: %s", line)
		}
		for _, h := range hexes {
			parsed, err := strconv.ParseInt(h, 16, 9)
			if err != nil {
				return nil, err
			}
			result.WriteByte(byte(parsed))
		}
	}
	return result.Bytes(), nil
}
