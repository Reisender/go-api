package api

type ParseError struct {
	Raw []byte
	Err error
}

func (pe ParseError) Error() string {
	return pe.Err.Error()
}
