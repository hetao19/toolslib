package slice

func UniqueInt64(s []int64) []int64 {
	size := len(s)
	if size == 0 {
		return []int64{}
	}
	m := make(map[int64]bool)
	for i:=0;i<size;i++ {
		m[s[i]] = true
	}

	_result := make([]int64, len(m))
	cnt := 0
	for k := range m {
		_result[cnt] = k
		cnt++
	}
	return _result
}

func UniqueInt(s []int) []int{
	size := len(s)
	if size == 0 {
		return []int{}
	}
	m := make(map[int]bool)
	for i:=0;i<size;i++ {
		m[s[i]] = true
	}

	_result := make([]int, len(m))
	cnt := 0
	for k := range m {
		_result[cnt] = k
		cnt++
	}
	return _result
}

func UniqueString(s []string) []string {
	size := len(s)
	if size == 0 {
		return []string{}
	}
	m := make(map[string]bool)
	for i:=0;i<size;i++ {
		m[s[i]] = true
	}
	_result := make([]int, len(m))
	cnt := 0
	for k := range m {
		_result[cnt] = k
		cnt++
	}
	return _result
}