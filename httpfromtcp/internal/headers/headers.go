package headers

import (
	"bytes"
	"fmt"
	"unicode"
)

var CRLF []byte = []byte("\r\n\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	consumed := 0

	if !bytes.Contains(data, CRLF) {
		return 0, false, nil
	}

	if bytes.HasPrefix(data, CRLF) {
		return 0, true, nil
	}

	splitData := bytes.Split(data, CRLF)

	headerUnParsed := splitData[0]

	consumed += len(headerUnParsed) + 4

	headerUnParsed = bytes.TrimSpace(headerUnParsed)

	headerBody := bytes.SplitN(headerUnParsed, []byte(" "), 2)

	fieldName := headerBody[0]
	fieldValue := headerBody[1]

	if err = checkValidFieldName(fieldName); err != nil {
		return 0, false, err
	}

	if _, exist := h[string(fieldName[:len(fieldName)-1])]; exist {
		h[string(fieldName[:len(fieldName)-1])] += ", " + string(fieldValue)
	} else {
		h[string(fieldName[:len(fieldName)-1])] = string(fieldValue)
	}

	return consumed, false, nil
}

func checkValidFieldName(fieldName []byte) error {

	validCharacter := []byte("!#$%^&'*-*._`|~")

	for i, char := range fieldName {
		if i == len(fieldName)-1 && char != byte(':') {
			return fmt.Errorf("invalid structure in field-line missing ':'")
		}
		if i == len(fieldName)-1 && char == byte(':') {
			continue
		}
		if unicode.IsLetter(rune(char)) || unicode.IsDigit(rune(char)) {
			continue
		}
		if bytes.Contains(validCharacter, []byte{char}) {
			continue
		}

		return fmt.Errorf("invalid character in field-line")
	}
	return nil
}
