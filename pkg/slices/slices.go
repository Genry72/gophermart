package slices

// Inslice проверяет наличие элемента в списке
func Inslice[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
