// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package rsrc_test

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/st3fan/rsrc"
	"github.com/stretchr/testify/assert"
)

func Test_Open(t *testing.T) {
	file, err := rsrc.Open("testdata/SolarianII")
	assert.NotNil(t, file)
	assert.Nil(t, err)

	assert.Equal(t, file.CountResources("CODE"), 5)

	for i := 0; i < file.CountResources("CODE"); i++ {
		_, ok := file.GetResource("CODE", i)
		assert.True(t, ok)
	}

	_, ok := file.GetResource("CODE", 6)
	assert.False(t, ok)

	assert.Equal(t, file.CountResources("snd "), 29)
	snd, ok := file.GetResource("snd ", 0)
	assert.True(t, ok)
	assert.Equal(t, snd.Name, "Present Bounce")
	assert.Equal(t, snd.ID, int16(5016))

	hex.Dumper(os.Stdout).Write(snd.Data)
}
