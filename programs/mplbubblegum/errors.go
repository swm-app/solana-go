package mplbubblegum

// fieldError and optionError report the exact field where a borsh encoding
// failure occurred, mirroring the error types the fork's generated packages
// (e.g. mplcore) produce so callers can rely on a consistent error shape
// across hand-written and generated packages.

type fieldError struct {
	Field string
	Err   error
}

func (e *fieldError) Error() string {
	return e.Field + ": " + e.Err.Error()
}

func (e *fieldError) Unwrap() error { return e.Err }

func newFieldError(field string, err error) error {
	if err == nil {
		return nil
	}
	return &fieldError{Field: field, Err: err}
}

type optionError struct {
	Field string
	Err   error
}

func (e *optionError) Error() string {
	return "?" + e.Field + ": " + e.Err.Error()
}

func (e *optionError) Unwrap() error { return e.Err }

func newOptionError(field string, err error) error {
	if err == nil {
		return nil
	}
	return &optionError{Field: field, Err: err}
}
