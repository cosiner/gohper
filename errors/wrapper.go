package errors

type Wrapper func(error) error

func (wrap Wrapper) Wrap(err error) error {
	if wrap == nil {
		return err
	}

	return wrap(err)
}
