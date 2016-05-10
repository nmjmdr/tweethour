package tweethour

import (
	"fmt"
)

type ErrorType int

const (
	ErrorStatus          = 1
	TokenGeneraionFailed = 2
	TimelineFetchFailed  = 3
)

type Error interface {
	error
	Type() ErrorType
}

type StatusError struct {
	StatusCode int
	StatusText string
}

func (s *StatusError) Error() string {
	return fmt.Sprintf("Request to get user's timeline failed, Status : %s", s.StatusText)
}

func (s *StatusError) Type() ErrorType {
	return ErrorStatus
}

func NewStatusError(code int, status string) *StatusError {
	s := new(StatusError)
	s.StatusCode = code
	s.StatusText = status
	return s
}

type TimelineError struct {
	e error
}

func (t *TimelineError) Type() ErrorType {
	return TimelineFetchFailed
}

func (t *TimelineError) Error() string {
	return t.e.Error()
}

func NewTimelineError(e error) *TimelineError {
	te := new(TimelineError)
	te.e = e
	return te
}

type TokenError struct {
	e error
}

func (t *TokenError) Type() ErrorType {
	return TokenGeneraionFailed
}

func (t *TokenError) Error() string {
	return t.e.Error()
}

func NewTokenError(e error) *TokenError {
	te := new(TokenError)
	te.e = e
	return te
}
