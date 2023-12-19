package main

type ErrBadRequest struct{}

func (e ErrBadRequest) Error() string {
	return "Bad Request"
}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "Not Found"
}

type ErrInternalServerError struct{}

func (e ErrInternalServerError) Error() string {
	return "Internal Server Error"
}

type ErrJsonParsing struct {
	Message string
}

func (e ErrJsonParsing) Error() string {
	return e.Message
}
