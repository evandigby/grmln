package grmln

import "fmt"

type responseError struct {
	response Response
}

func (e responseError) Response() Response {
	return e.response
}

func (e responseError) Error() string {
	return fmt.Sprintf("Error (%d: %q) in request %q: %q", e.response.Status.Code, StatusString(e.response.Status.Code), e.response.RequestID, e.response.Status.Message)
}

func (e responseError) IsUnauthorized() bool {
	return e.response.Status.Code == StatusUnauthorized
}
func (e responseError) IsAuthenticate() bool {
	return e.response.Status.Code == StatusAuthenticate
}
func (e responseError) IsMalformedRequest() bool {
	return e.response.Status.Code == StatusMalformedRequest
}
func (e responseError) IsInvalidRequestArguments() bool {
	return e.response.Status.Code == StatusInvalidRequestArguments
}
func (e responseError) IsServerError() bool {
	return e.response.Status.Code == StatusServerError
}
func (e responseError) IsScriptEvaluationError() bool {
	return e.response.Status.Code == StatusScriptEvaluationError
}
func (e responseError) IsServerTimeout() bool {
	return e.response.Status.Code == StatusServerTimeout
}
func (e responseError) IsServerSerializationError() bool {
	return e.response.Status.Code == StatusServerSerializationError
}
func (e responseError) IsInvalidResponseCode() bool {
	return e.response.Status.Code.IsInvalid()
}

type unauthorized interface {
	IsUnauthorized() bool
}
type authenticate interface {
	IsAuthenticate() bool
}
type malformedRequest interface {
	IsMalformedRequest() bool
}
type invalidRequestArguments interface {
	IsInvalidRequestArguments() bool
}
type serverError interface {
	IsServerError() bool
}
type scriptEvaluationError interface {
	IsScriptEvaluationError() bool
}
type serverTimeout interface {
	IsServerTimeout() bool
}
type serverSerializationError interface {
	IsServerSerializationError() bool
}
type invalidResponseCode interface {
	IsInvalidResponseCode() bool
}

// IsUnauthorized returns whether or not the error is a Unauthorized error
func IsUnauthorized(err error) bool {
	e, ok := err.(unauthorized)
	return ok && e.IsUnauthorized()
}

// IsAuthenticate returns whether or not the error is a Authenticate error
func IsAuthenticate(err error) bool {
	e, ok := err.(authenticate)
	return ok && e.IsAuthenticate()
}

// IsMalformedRequest returns whether or not the error is a MalformedRequest error
func IsMalformedRequest(err error) bool {
	e, ok := err.(malformedRequest)
	return ok && e.IsMalformedRequest()
}

// IsInvalidRequestArguments returns whether or not the error is a InvalidRequestArguments error
func IsInvalidRequestArguments(err error) bool {
	e, ok := err.(invalidRequestArguments)
	return ok && e.IsInvalidRequestArguments()
}

// IsServerError returns whether or not the error is a ServerError error
func IsServerError(err error) bool {
	e, ok := err.(serverError)
	return ok && e.IsServerError()
}

// IsScriptEvaluationError returns whether or not the error is a ScriptEvaluationError error
func IsScriptEvaluationError(err error) bool {
	e, ok := err.(scriptEvaluationError)
	return ok && e.IsScriptEvaluationError()
}

// IsServerTimeout returns whether or not the error is a ServerTimeout error
func IsServerTimeout(err error) bool {
	e, ok := err.(serverTimeout)
	return ok && e.IsServerTimeout()
}

// IsServerSerializationError returns whether or not the error is a ServerSerializationError error
func IsServerSerializationError(err error) bool {
	e, ok := err.(serverSerializationError)
	return ok && e.IsServerSerializationError()
}

// IsInvalidResponseCode returns whether or not the response code is completely invalid (not in the supported spec)
func IsInvalidResponseCode(err error) bool {
	e, ok := err.(invalidResponseCode)
	return ok && e.IsInvalidResponseCode()
}

