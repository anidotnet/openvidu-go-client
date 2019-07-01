package openvidu

import "fmt"

type openViduError struct {
	Status int	`json:"status"`
}

func newOpenViduError(status int) *openViduError {
	return &openViduError{Status: status}
}

func (err *openViduError) Error() string {
	return fmt.Sprintf("Invalid status code %d recieved from OpenVidu server", err.Status)
}
