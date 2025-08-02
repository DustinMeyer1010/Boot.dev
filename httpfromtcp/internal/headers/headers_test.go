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
	assert.Equal(t, 23, n)
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
	assert.Equal(t, 30, n)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069")

	// End of Header Should return done as true
	headers = NewHeaders()
	data = []byte("Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["remain"] = "This should not be removed"
	data = []byte("Host: localhost:42069\r\nType: String\r\n\r\n")
	n, done, err = headers.Parse(data)
	data = data[n:]
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, n, 23)
	assert.Equal(t, headers["Host"], "localhost:42069")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 14)
	assert.Equal(t, headers["Type"], "String")
	assert.False(t, done)

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
	assert.Equal(t, n, 25)
	assert.False(t, done)

	// Valid character in feild-name
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:3000\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 23)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069")

	data = data[n:]

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, n, 22)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069, localhost:3000")

}
