package errors

type Wrapper func(error) error

func (wrap Wrapper) Wrap(err error) error {
	if wrap == nil {
		return err
	}

	return wrap(err)
}

type WrappedError interface {
	Unwrap() error
}

func Unwrap(err error) error {
	for {
		if err == nil {
			return nil
		}

		e, is := err.(WrappedError)
		if is {
			err = e.Unwrap()
		} else {
			break
		}
	}

	return err
}
