package spatial

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWKBNullLineString(t *testing.T) {
	b, _ := hex.DecodeString("010300000030303000000000003030303030303030")
	buf := bytes.NewBuffer(b)

	var g Geom
	err := g.UnmarshalWKB(buf)
	assert.NotNil(t, err)
}
