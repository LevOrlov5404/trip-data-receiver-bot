package models

type ClientError struct {
	Description string
}

func (clientError ClientError) Error() string {
	return clientError.Description
}
