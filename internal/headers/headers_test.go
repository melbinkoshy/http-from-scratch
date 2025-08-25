package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid  header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nfoo: fighters\r\nfoo:  are the best        \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, _ = headers.Get("Host")
	foo, _ := headers.Get("foo")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "fighters, are the best", foo)
	assert.Equal(t, 68, n)
	assert.True(t, done)

	// Test: Valid  header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nfoo: fighters        \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, _ = headers.Get("Host")
	foo, _ = headers.Get("foo")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "fighters", foo)
	assert.Equal(t, 48, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in key
	headers = NewHeaders()
	data = []byte("      H©st: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in key
	headers = NewHeaders()
	data = []byte("      Ho st: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
