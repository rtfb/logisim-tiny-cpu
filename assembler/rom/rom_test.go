package rom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseROMDump(t *testing.T) {
	input := `v3.0 hex words addressed
00: 0f 10 0e 18 49 08 24 29 08 54 34 08 c0 00 00 00
10: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00`
	if !strings.HasPrefix(input, "v3.0 hex words addressed") {
		t.Fatal("only v3.0 hex is supported")
	}
	b, err := parse(input)
	assert.NoError(t, err)
	assert.Len(t, b, 32)
	assert.Equal(t, []byte{0x0f, 0x10, 0x0e, 0x18, 0x49}, b[:5])
}
