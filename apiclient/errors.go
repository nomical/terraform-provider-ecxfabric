package apiclient

import "fmt"

type oAuthErrorResponse struct {
	ErrorDomain      string `json:"errorDomain"`      // example: apps-fqa
	ErrorTitle       string `json:"errorTitle"`       // example: Invalid Username/Password
	ErrorCode        string `json:"errorCode"`        // example: S1003
	DeveloperMessage string `json:"developerMessage"` // example: Invalid Username/Password
	ErrorMessage     string `json:"errorMessage"`     // example: We didn't recognize the username or password you entered. Please try again..
}

type errorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	MoreInfo     string `json:"moreInfo"`
	Property     string `json:"property"`
}

type errorResponseArray []errorResponse

type AuthError struct {
	HTTPStatusCode int
	oAuthErrorResponse
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("HTTPStatusCode: %v, ErrorDomain: %v, ErrorTitle: %v, ErrorCode: %v, DeveloperMessage: %v, ErrorMessage: %v", e.HTTPStatusCode, e.ErrorDomain, e.ErrorTitle, e.ErrorCode, e.DeveloperMessage, e.ErrorMessage)
}

type Error struct {
	HTTPStatusCode int
	errorResponse
}

func (e *Error) Error() string {
	return fmt.Sprintf("HTTPStatusCode: %v, ErrorCode: %v, ErrorMessage: %v, MoreInfo: %v, Property: %v", e.HTTPStatusCode, e.ErrorCode, e.ErrorMessage, e.MoreInfo, e.Property)
}

type ErrorArray struct {
	HTTPStatusCode int
	errorResponseArray
}

func (ea *ErrorArray) Error() string {
	errStr := fmt.Sprintf("HTTPStatusCode: %v\r\n", ea.HTTPStatusCode)
	for _, e := range ea.errorResponseArray {
		errStr += fmt.Sprintf("ErrorCode: %v, ErrorMessage: %v, MoreInfo: %v, Property: %v\r\n", e.ErrorCode, e.ErrorMessage, e.MoreInfo, e.Property)
	}

	return errStr
}
