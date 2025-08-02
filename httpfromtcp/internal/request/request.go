package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/DustinMeyer1010/httpfromtcp/internal/headers"
)

type Status int

const (
	initialized Status = iota
	requestStateParsingHeaders
	done
)

type Request struct {
	Status      Status
	Headers     headers.Headers
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var buf = make([]byte, 8)

	request := Request{Status: initialized, Headers: headers.NewHeaders()}

	readToIndex := 0

	for request.Status != done {

		if readToIndex >= len(buf) {
			if len(buf) == 0 {
				buf = make([]byte, 8)
			} else {
				tempbuf := make([]byte, len(buf)*2)
				copy(tempbuf, buf)
				buf = tempbuf
			}

		}

		n, err := reader.Read(buf[readToIndex:])

		if err == io.EOF || err != nil {
			break
		}

		readToIndex += n

		consumed, err := request.parse(buf)

		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			buf = buf[consumed:]
			readToIndex -= consumed
		}

	}

	fmt.Println("Before Return ", request.Headers)

	return &request, nil
}

func parseRequestLine(request []byte) (int, error) {

	requestParse := bytes.Split(request, []byte("\r\n"))

	if len(requestParse) == 1 {
		return 0, nil
	}

	return len(requestParse[0]) + 2, nil

}

func (r *Request) parse(data []byte) (int, error) {

	var consumed int
	var err error
	var complete bool

	fmt.Printf("%q\n", data)

	switch r.Status {
	case done:
		return 0, fmt.Errorf("trying to read data after done")
	case initialized:
		consumed, err = parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if consumed == 0 {
			return 0, nil
		}

		requestLine := bytes.Split(data[:consumed-2], []byte(" "))

		if len(requestLine) != 3 {
			return 0, fmt.Errorf("invalid request line")
		}

		method := string(requestLine[0])
		path := string(requestLine[1])
		version := string(requestLine[2][5:])

		if err := checkValidMethod(method); err != nil {
			return 0, err
		}

		if err := checkValidVersion(version); err != nil {
			return 0, err
		}

		r.Status = requestStateParsingHeaders

		r.RequestLine = RequestLine{
			Method:        method,
			RequestTarget: path,
			HttpVersion:   version,
		}
	case requestStateParsingHeaders:
		consumed, complete, err = r.Headers.Parse(data)

		if err != nil {
			return 0, err
		}

		if complete {
			r.Status = done
		}
	}
	return consumed, nil
}

func checkValidMethod(method string) error {
	switch method {
	case "POST", "GET", "OPTION", "HEAD", "PUT", "DELETE", "TRACE", "PATCH", "CONNECT":
		return nil
	default:
		return fmt.Errorf("invalid method")
	}
}

func checkValidVersion(version string) error {
	if version != "1.1" {
		return fmt.Errorf("invalid version")
	}

	return nil
}
