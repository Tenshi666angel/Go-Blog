package servererror

import "errors"

var (
    InvalidJson = errors.New("invalid json from client")
    InternalError = errors.New("internal server error")
    InvalidCrerdentials = errors.New("invalid username or password")
    ResourceNotFound = errors.New("resource not found")
)
