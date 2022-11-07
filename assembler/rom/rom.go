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

func Dump(b []byte) string {
	var buf strings.Builder
	buf.WriteString("v3.0 hex words addressed\n")
	i := 0
	for i < 256 {
		if i%16 == 0 {
			if i > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(fmt.Sprintf("%02x:", i))
		}
		var out byte
		if i < len(b) {
			out = b[i]
		}
		buf.WriteString(fmt.Sprintf(" %02x", out))
		i++
	}
	return buf.String() + "\n"
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
