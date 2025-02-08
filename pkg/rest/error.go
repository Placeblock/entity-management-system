package rest

import "fmt"

type HTTPError struct {
	Response
	Cause error `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("title: %s, detail: %s", e.Title, e.Detail)
	}
	return fmt.Sprintf("title: %s, detail: %s, cause: %s", e.Title, e.Detail, e.Cause)
}
