package w

func M[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Error[T any](_ T, err error) error {
	return err
}

func ToInterfaces[T any](objs []T) []interface{} {
	var newObjs []interface{}
	for _, obj := range objs {
		newObjs = append(newObjs, obj)
	}
	return newObjs
}

func P[T any](o T) *T {
	return &o
}
