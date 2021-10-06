package investing

type ParsingError struct {
	Html string
	Err  error
}

func (e *ParsingError) Error() string {
	return e.Err.Error()
}
