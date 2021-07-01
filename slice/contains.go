package slice

func Contains(s1 []interface{}, v interface{}) bool {
	for _, value := range s1 {
		if value == v {
			return true
		}
	}
	return false
}

func ContainsInt(s1 []int, v int) bool {
	for _, value := range s1 {
		if value == v {
			return true
		}
	}
	return false
}

func ContainsInt64(sl []int64, v int64) bool {
	for _, value := range sl {
		if value == v {
			return true
		}
	}
	return false
}

func ContainsString(sl []string, v string) bool {
	for _, value := range sl {
		if value == v {
			return true
		}
	}
	return false
}
