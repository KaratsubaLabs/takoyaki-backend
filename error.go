package main

type HTTPError interface {
	error
	Status() int
}

type HTTPStatusError inteface {
	Code        int
	Err         error
}

func (se HTTPStatusError) Error() string {
	return se.Err.Error()
}

func (set HTTPStatusError) Status() int {
	return se.Code
}
