package utils

import "errors"

func ConvertSliceTo[T any](iter []any) (sol []T, err error) {
	sol = make([]T, len(iter))

	for i := range iter {
		conv, ok := iter[i].(T)
		if !ok {
			return nil, errors.New("different types")
		}
		sol[i] = conv
	}
	return sol, nil
}
