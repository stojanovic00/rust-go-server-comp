package main

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"
)

type HttpMethod int8

const (
	Get HttpMethod = iota
	Put
	NotImplemented
)

type HttpRequest struct {
	Method  HttpMethod
	Path    string
	Headers []string
	Body    *Entity
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{
		Method:  NotImplemented,
		Path:    "",
		Headers: make([]string, 0),
		Body:    nil,
	}
}

func ParseRequest(conn net.Conn) (*HttpRequest, error) {
	reader := bufio.NewReader(conn)

	request := NewHttpRequest()

	//Parse headers
	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrBadRequest{}
		}

		if header == "\r\n" {
			break
		}

		request.Headers = append(request.Headers, strings.TrimSpace(header))
	}

	//Parse request line
	request_line := request.Headers[0]
	request.Headers = request.Headers[1:]

	request_parts := strings.Split(request_line, " ")
	for idx, _ := range request_parts {
		request_parts[idx] = strings.TrimSpace(request_parts[idx])
	}

	switch request_parts[0] {
	case "GET":
		request.Method = Get
	case "PUT":
		request.Method = Put
	default:
		request.Method = NotImplemented
	}

	request.Path = request_parts[1]

	//Parse body
	contentLength := 0
	for _, header := range request.Headers {
		if strings.Contains(header, "Content-Length") {
			parts := strings.Split(header, ":")
			contentLengthStr := strings.TrimSpace(parts[1])
			var err error
			contentLength, err = strconv.Atoi(contentLengthStr)
			if err != nil {
				return nil, ErrBadRequest{}
			}
		}
	}

	if contentLength > 0 {
		bodyBytes := make([]byte, contentLength)
		readNum, err := reader.Read(bodyBytes)
		if err != nil || readNum != contentLength {
			return nil, ErrBadRequest{}
		}

		var entity Entity
		err = json.Unmarshal([]byte(bodyBytes), &entity)
		if err != nil {
			return nil, ErrJsonParsing{Message: err.Error()}
		}

		request.Body = &entity
	}

	return request, nil
}
