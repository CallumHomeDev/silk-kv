package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// ParseRequest parse a request from a reader and return a Request object
func ParseRequest(reader *bufio.Reader) ([][]byte, error) {
	// Read first line, ex: *2\r\n
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 1 || line[0] != '*' {
		return nil, fmt.Errorf("expected '*', got %q", line)
	}

	// Convert to number
	count, err := strconv.Atoi(string(bytes.TrimSpace(line[1:])))
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %v", err)
	}

	args := make([][]byte, count)
	for i:= 0; i < count; i++ {
		// Read next line, ex: $3\r\nabc\r\n
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		if len(line) < 1 || line[0] != '$' {
			return nil, fmt.Errorf("expected '$', got %q", line)
		}
		length, err := strconv.Atoi(string(bytes.TrimSpace(line[1:])))
		if err != nil {
			return nil, fmt.Errorf("invalid bulk length: %v", err)
		}

		// Read paytload (length bytes) + \r\n
		buf := make([]byte, length+2)
		if _, err = io.ReadFull(reader, buf); err != nil {
			return nil, err
		}
		args[i] = buf[:length]
	}

	return args, nil
}

// ParseFrame parse a RESP frame from a reader and return the frame type and payload
func ParseFrame(data []byte) (args [][]byte, consumed int, err error) {
	// Wrap the data in a buffer
	r := bytes.NewReader(data)
	reader := bufio.NewReader(r)

	// Call ParseRequest to parse the request
	args, err = ParseRequest(reader)
	if err != nil {
		return nil, 0, err
	}

	// Count the number of bytes consumed
	consumed = len(data) - reader.Buffered()
	return args, consumed, nil
}


// FormatSimpleString response RESP simple string: +<message>\r\n
func FormatSimpleString(message string) []byte {
	return []byte("+" + message + "\r\n")
}

// FormatError response RESP error: -<message>\r\n
func FormatError(message string) []byte {
    return []byte("-" + message + "\r\n")
}

// FormatBulkString response RESP bulk string: $<len>\r\n<payload>\r\n
func FormatBulkString(payload []byte) []byte {
    return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(payload), payload))
}

// FormatNilBulkString response RESP nil bulk string: $-1\r\n
func FormatNilBulkString() []byte {
    return []byte("$-1\r\n")
}

// FormatInteger response RESP integer: :<number>\r\n
func FormatInteger(n int64) []byte {
	return []byte(":" + strconv.FormatInt(n, 10) + "\r\n")
}

// FormatArray response RESP array: *<count>\r\n
func FormatArray(arr [][]byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(len(arr)) + "\r\n")
	for _, elem := range arr {
		buf.WriteString("$" + strconv.Itoa(len(elem)) + "\r\n")
		buf.Write(elem)
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}