package serde

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Sanity check on marshal / unmarshal
func TestJSON(t *testing.T) {
	serde := JSON{}

	data := map[string]string{
		"a": "b",
		"1": "2",
	}

	buff, err := serde.Marshal(&data)
	require.NoError(t, err)

	result := map[string]string{}

	err = serde.Unmarshal(buff, &result)
	require.NoError(t, err)

	require.Equal(t, data, result)
}
