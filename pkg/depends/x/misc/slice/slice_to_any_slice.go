package slice

func ToAnySlice[K any](vs ...K) []interface{} {
	interfaces := make([]any, 0, len(vs))
	for _, v := range vs {
		interfaces = append(interfaces, v)
	}
	return interfaces
}
