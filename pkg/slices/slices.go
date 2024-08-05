package slices

func Map[T, E any](list []T, f func(T) E) []E {
	newList := make([]E, 0, len(list))
	for _, i := range list {
		newList = append(newList, f(i))
	}
	return newList
}
