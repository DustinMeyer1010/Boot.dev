package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {

	//Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Header with extra spaces
	headers = NewHeaders()
	data = []byte("Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 32, n)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069")

	// End of Header Should return done as true
	headers = NewHeaders()
	data = []byte("\r\n\r\nHost : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["remain"] = "This should not be removed"
	data = []byte("Host: localhost:42069\r\n\r\nType: String\r\n\r\n")
	n, done, err = headers.Parse(data)
	data = data[n:]
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, n, 25)
	assert.Equal(t, headers["Host"], "localhost:42069")
	n, done, err = headers.Parse(data)
	data = data[n:]
	require.NoError(t, err)
	assert.Equal(t, n, 16)
	assert.Equal(t, headers["Type"], "String")
	assert.False(t, done)
	assert.Equal(t, len(data), 0)

	// Invalid field-line with special character
	headers = NewHeaders()
	data = []byte("H@ost: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, n, 0)
	assert.False(t, done)

	// Valid character in feild-name
	headers = NewHeaders()
	data = []byte("Ho1st!: localhost:42069\r\n\r\n")
	fmt.Println(len(data))
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 27)
	assert.False(t, done)

	// Valid character in feild-name
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\nHost: localhost:3000\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 25)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069")

	data = data[n:]

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 24)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069, localhost:3000")

}
