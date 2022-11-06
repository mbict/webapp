package container

var DefaultContainer Container

// Get tries to resolve the type T from the container and returns the instance.
// If it cannot resolve or one of the containers returns an error it will panic
func Get[T any](container ...Container) T {
	val, err := GetE[T](container...)
	if err != nil {
		panic(err)
	}
	return val
}

func GetE[T any](container ...Container) (val T, err error) {
	if len(container) == 0 {
		container = []Container{DefaultContainer}
	}

	for _, c := range container {
		err = c.InvokeE(func(a T) {
			val = a
		})

		if err == nil {
			return val, nil
		}

		if err != ErrCannotResolveInstance {
			return val, err
		}
	}
	return val, ErrCannotResolveInstance
}
