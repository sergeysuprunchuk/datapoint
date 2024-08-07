package slices

func Map[T, E any](list []T, f func(T) E) []E {
	newList := make([]E, 0, len(list))
	for _, i := range list {
		newList = append(newList, f(i))
	}
	return newList
}

func Unique[T comparable](s ...[]T) []T {
	result, dict := make([]T, 0), make(map[T]struct{})
	for _, i := range s {
		for _, j := range i {
			if _, ok := dict[j]; !ok {
				result = append(result, j)
			}
		}
	}
	return result
}
