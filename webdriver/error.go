package webdriver

import (
	"errors"
)

type ErrorCode string

var (
	ElementClickIntercepted ErrorCode = "element click intercepted"
	ElementNotInteractable  ErrorCode = "element not interactable"
	InsecureCertificate     ErrorCode = "insecure certificate"
	InvalidArgument         ErrorCode = "invalid argument"
	InvalidCookieDomain     ErrorCode = "invalid cookie domain"
	InvalidElementState     ErrorCode = "invalid element state"
	InvalidSelector         ErrorCode = "invalid selector"
	InvalidSessionID        ErrorCode = "invalid session id"
	JavascriptError         ErrorCode = "javascript error"
	MoveTargetOutOfBounds   ErrorCode = "move target out of bounds"
	NoSuchAlert             ErrorCode = "no such alert"
	NoSuchCookie            ErrorCode = "no such cookie"
	NoSuchElement           ErrorCode = "no such element"
	NoSuchFrame             ErrorCode = "no such frame"
	NoSuchWindow            ErrorCode = "no such window"
	ScriptTimeout           ErrorCode = "script timeout"
	SessionNotCreated       ErrorCode = "session not created"
	StaleElementReference   ErrorCode = "stale element reference"
	Timeout                 ErrorCode = "timeout"
	UnableToSetCookie       ErrorCode = "unable to set cookie"
	UnableToCaptureScreen   ErrorCode = "unable to capture screen"
	UnexpectedAlertOpen     ErrorCode = "unexpected alert open"
	UnknownCommand          ErrorCode = "unknown command"
	UnknownError            ErrorCode = "unknown error"
	UnknownMethod           ErrorCode = "unknown method"
	UnsupportedOperation    ErrorCode = "unsupported operation"
)

type Error struct {
	ErrorCode  ErrorCode
	Message    string
	StackTrace []string
}

func (e *Error) Error() string {
	return "webdriver: " + string(e.ErrorCode) + ": " + e.Message
}

func IsNoSuchElement(err error) bool {
	var wErr *Error
	return errors.As(err, &wErr) && wErr.ErrorCode == NoSuchElement
}
