package domain

// BadRequestError represents an error generated by a bad request
type BadRequestError struct{ Err error }

func (b BadRequestError) Error() string {
	return b.Err.Error()
}

// ValidationError represents an error that refers to a failed validation
type ValidationError struct{ Err error }

func (v ValidationError) Error() string {
	return v.Err.Error()
}
