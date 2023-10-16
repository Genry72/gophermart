package pkg

// ToSliceInterface приводит массив к массиву интерфейсов
func ToSliceInterface[T comparable](elems []T) []interface{} {
	result := make([]interface{}, len(elems))

	for i, v := range elems {
		result[i] = v
	}

	return result
}
