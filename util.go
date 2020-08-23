package rsrc

import (
	"bufio"
	"io"
)

func fourCharacterCode(code uint32) string {
	bytes := []byte{
		byte(code >> 24 & 0x000000ff),
		byte(code >> 16 & 0x000000ff),
		byte(code >> 8 & 0x000000ff),
		byte(code & 0x000000ff),
	}

	return string(bytes)
}

func readPascalString(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)

	length, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	name := ""
	for i := 0; i < int(length); i++ {
		c, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		name += string(c)
	}

	return name, nil
}
