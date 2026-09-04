package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type parseState string

type Request struct {
	RequestLine  RequestLine
	RequestState parseState
}

func NewRequest() *Request {
	return &Request{
		RequestState: StateInit,
	}
}

const (
	StateInit  parseState = "init"
	StateError parseState = "error"
	StateDone  parseState = "done"
)

var (
	ErrBadReqLine = errors.New("bad request line")
	SEPERATOR     = []byte("\r\n")
	ErrReq        = errors.New("Request in error")
)

func parseRequestLine(request []byte) (*RequestLine, int, error) {

	idx := bytes.Index(request, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	reqLine := request[:idx]

	readN := idx + len(SEPERATOR)

	reqLineParts := bytes.Split(reqLine, []byte(" "))

	if len(reqLineParts) != 3 {
		return nil, 0, fmt.Errorf("Improper request-line")
	}
	httpVersionParts := bytes.Split(reqLineParts[2], []byte("/"))

	if len(httpVersionParts) != 2 || string(httpVersionParts[0]) != "HTTP" || string(httpVersionParts[1]) != "1.1" {
		return nil, 0, fmt.Errorf("Improper Http versioning")
	}

	return &RequestLine{
		string(reqLineParts[0]),
		string(reqLineParts[1]),
		string(httpVersionParts[1]),
	}, readN, nil

}

func (r *Request) parse(data []byte) (int, error) {

	readN := 0
outer:
	for {
		switch r.RequestState {
		case StateError:
			return 0, ErrReq
		case StateInit:
			rl, n, err := parseRequestLine(data)
			if err != nil {
				r.RequestState = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			readN += n

			r.RequestState = StateDone

		case StateDone:
			break outer

		}
	}
	return readN, nil
}

func (r *Request) done() bool {
	return r.RequestState == StateDone || r.RequestState == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := NewRequest()

	buffer := make([]byte, 1024)
	bufferIdx := 0

	for !request.done() {
		n, err := reader.Read(buffer[bufferIdx:])
		if err != nil {
			return nil, err
		}

		bufferIdx += n

		readN, err := request.parse(buffer[:bufferIdx])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[readN:bufferIdx])
		bufferIdx -= readN

	}

	return request, nil
}
