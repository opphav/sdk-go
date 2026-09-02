package services

func (i *Service) Execute(d Dimension, fn func() error) (err error) {
	i.Start(d)
	defer i.Stop(d)

	defer func() {
		if r := recover(); r != nil {
			i.Fail(d)
			panic(r)
		} else if err != nil {
			i.Fail(d)
		} else {
			i.Success(d)
		}
	}()

	return fn()
}

func (i *Service) Dispatch(d Dimension, fn func() error) (err error) {
	i.Call(d)

	defer func() {
		if r := recover(); r != nil {
			i.Fail(d)
			panic(r)
		} else if err != nil {
			i.Fail(d)
		} else {
			i.Success(d)
		}
	}()

	return fn()
}
