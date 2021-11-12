package main

type HTTPError interface {
	error
	Status() int
}

type HTTPStatusError struct {
	Code        int
	Err         error
}

func (se HTTPStatusError) Error() string {
	return se.Err.Error()
}

func (se HTTPStatusError) Status() int {
	return se.Code
}
