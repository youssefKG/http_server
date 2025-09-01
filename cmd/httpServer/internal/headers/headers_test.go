// Test: Valid single header
package headers

// func TestHeaders(t *testing.T) {
// 	headers := NewHeaders()
// 	data := []byte("Host: localhost:42069\r\n\r\n")
// 	n, done, err := headers.Parse(data)
// 	require.NoError(t, err)
// 	require.NotNil(t, headers)
// 	assert.Equal(t, "localhost:42069", headers["Host"])
// 	fmt.Println("n: ", n)
// 	assert.False(t, done)
//
// 	// Test: Invalid spacing header
// 	headers = NewHeaders()
// 	data = []byte("       Host : localhost:42069       \r\n\r\n")
// 	n, done, err = headers.Parse(data)
// 	require.Error(t, err)
// 	assert.Equal(t, 0, n)
// 	assert.False(t, done)
//
// }
