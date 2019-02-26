package grmln

import (
	"encoding/json"
	"fmt"
)

// Response contains a response message
type Response struct {
	RequestID string         `json:"requestId"`
	Status    ResponseStatus `json:"status"`
	Result    ResponseResult `json:"result"`
}

// ResponseStatus contains the status details about the response
type ResponseStatus struct {
	Code       StatusCode             `json:"code"`
	Attributes map[string]interface{} `json:"attributes"`
	Message    string                 `json:"message"`
}

// ResponseResult contains the result of the response
type ResponseResult struct {
	Data json.RawMessage        `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

// IsPartial returns whether or not the response has partial content
func (r Response) IsPartial() bool {
	return r.Status.Code == StatusPartialContent
}

// Err returns the response's error
func (r Response) Err() error {
	switch r.Status.Code {
	case StatusSuccess,
		StatusNoContent,
		StatusPartialContent:

		return nil

	// All other codes are error codes (including invalid ones)
	default:
		return responseError{response: r}
	}
}

// StatusCode is a type for gremlin status codes
type StatusCode int

// Status codes
const (
	// StatusSuccess means the server successfully processed a request to completion - there are no messages remaining in this stream. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusSuccess StatusCode = 200
	// StatusNoContent means the server processed the request but there is no result to return (e.g. an Iterator with no elements) - there are no messages remaining in this stream. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusNoContent StatusCode = 204
	// StatusPartialContent means the server successfully returned some content, but there is more in the stream to arrive - wait for a SUCCESS to signify the end of the stream. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusPartialContent StatusCode = 206
	// StatusUnauthorized means the request attempted to access resources that the requesting user did not have access to. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusUnauthorized StatusCode = 401
	// StatusAuthenticate means a challenge from the server for the client to authenticate its request. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusAuthenticate StatusCode = 407
	// StatusMalformedRequest means the request message was not properly formatted which means it could not be parsed at all or the "op" code was not recognized such that Gremlin Server could properly route it for processing. Check the message format and retry the request. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusMalformedRequest StatusCode = 498
	// StatusInvalidRequestArguments means the request message was parseable, but the arguments supplied in the message were in conflict or incomplete. Check the message format and retry the request. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusInvalidRequestArguments StatusCode = 499
	// StatusServerError means a general server error occurred that prevented the request from being processed. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusServerError StatusCode = 500
	// StatusScriptEvaluationError means the script submitted for processing evaluated in the ScriptEngine with errors and could not be processed. Check the script submitted for syntax errors or other problems and then resubmit. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusScriptEvaluationError StatusCode = 597
	// StatusServerTimeout means the server exceeded one of the timeout settings for the request and could therefore only partially responded or did not respond at all. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusServerTimeout StatusCode = 598
	// StatusServerSerializationError means the server was not capable of serializing an object that was returned from the script supplied on the request. Either transform the object into something Gremlin Server can process within the script or install mapper serialization classes to Gremlin Server. (see http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements)
	StatusServerSerializationError StatusCode = 599
)

// IsInvalid returns whether or not the status code is invalid
func (c StatusCode) IsInvalid() bool {
	switch c {
	case StatusSuccess,
		StatusNoContent,
		StatusPartialContent,
		StatusUnauthorized,
		StatusAuthenticate,
		StatusMalformedRequest,
		StatusInvalidRequestArguments,
		StatusServerError,
		StatusScriptEvaluationError,
		StatusServerTimeout,
		StatusServerSerializationError:

		return false
	default:
		return true
	}
}

var statusCodeStrings = map[StatusCode]string{
	StatusSuccess:                  "Success",
	StatusNoContent:                "No Content",
	StatusPartialContent:           "Partial Content",
	StatusUnauthorized:             "Unauthorized",
	StatusAuthenticate:             "Authenticate",
	StatusMalformedRequest:         "Malformed Request",
	StatusInvalidRequestArguments:  "Invalid Request Arguments",
	StatusServerError:              "Server Error",
	StatusScriptEvaluationError:    "Script Evaluation Error",
	StatusServerTimeout:            "Server Timeout",
	StatusServerSerializationError: "Server Serialization Error",
}

// StatusString returns the stringified version of the status code
func StatusString(code StatusCode) string {
	if code.IsInvalid() {
		return fmt.Sprintf("Invalid Response Code: %d", code)
	}

	return statusCodeStrings[code]
}

