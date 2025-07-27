package request

import (
	"bytes"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestParse := bytes.Split(request, []byte("\r\n"))
	requestLineAll := requestParse[0]

	requestLine, err := ParseRequestLine(requestLineAll)

	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, err
}

func ParseRequestLine(requestLine []byte) (*RequestLine, error) {

	requestLineList := bytes.Split(requestLine, []byte(" "))

	if len(requestLineList) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method := string(requestLineList[0])

	if err := checkValidMethod(method); err != nil {
		return nil, err
	}

	target := string(requestLineList[1])

	versionNumber := string(bytes.Split(requestLineList[2], []byte("/"))[1])

	if err := checkValidVersion(versionNumber); err != nil {
		return nil, err
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   versionNumber,
	}, nil
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
