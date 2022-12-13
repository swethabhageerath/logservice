package exceptions

type Exception struct {
	Err        error
	StatusCode int
}

func (e Exception) Error() string {
	return e.Err.Error()
}s
