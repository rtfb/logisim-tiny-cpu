package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInput(t *testing.T) {
	const src = `
/* C-style comment
 *
 */

li 1
getacc r7 // inline comment
getacc r6
`
	toks := Tokenize(strings.NewReader(src), "test")
	toksStr := []string{}
	for _, tok := range toks {
		toksStr = append(toksStr, tok.Text)
	}
	assert.Equal(t, []string{"li", "1", "getacc", "r7", "getacc", "r6"}, toksStr)
}
