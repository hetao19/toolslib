package slice

func MergeInterface(x, y []interface{}) (result []interface{}) {
	result = append(x, y...)
	return result
}

func MergeInt(x, y []int) (result []int) {
	result = append(x, y...)
	return result
}

func MergeInt64(x, y []int) (result []int) {
	result = append(x, y...)
	return result
}

func MergeString(x, y []string) (result []string) {
	result = append(x, y...)
	return result
}