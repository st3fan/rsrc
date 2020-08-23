package rsrc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fourCharacterCode(t *testing.T) {
	assert.Equal(t, "APPL", fourCharacterCode(1095782476))
}

func Test_readPascalString_Empty(t *testing.T) {
	s, e := readPascalString(bytes.NewReader([]byte{0}))
	assert.NoError(t, e)
	assert.Empty(t, s)
}

func Test_readPascalString_Good(t *testing.T) {
	s, e := readPascalString(bytes.NewReader([]byte{3, 'H', 'i', '!'}))
	assert.NoError(t, e)
	assert.Equal(t, "Hi!", s)
}

func Test_readPascalString_TooShort(t *testing.T) {
	_, e := readPascalString(bytes.NewReader([]byte{3, 'H', 'i'}))
	assert.Error(t, e)
}
