package api

// getDigits returns an array of individual digits of a number
// useful for formatting numbers to strings
func getDigits(i int) []int {
	if i == 0 {
		return []int{0}
	}
	var ret []int
	for i > 0 {
		ret = append(ret, i%10)
		i /= 10
	}
	return reverseInts(ret)
}

// reverseInts returns a reversed copy of an array of integers
func reverseInts(ints []int) []int {
	var ret []int
	for i := len(ints) - 1; i >= 0; i-- {
		ret = append(ret, ints[i])
	}
	return ret
}

// parseInt parses a string into an integer
// it is a zero-dependant version of strconv.Atoi(s)
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, NewError(ErrCodeInternalError, "cannot parse empty string")
	}
	var ret int
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, NewError(ErrCodeInternalError, "cannot parse string:"+s)
		}
		ret += int(c-'0') * pow(10, len(s)-i-1)
	}
	return ret, nil
}

// pow returns x raised to the power of y
// it is a zero-dependant version of math.Pow(x, y)
func pow(x, y int) int {
	ret := 1
	for i := 0; i < y; i++ {
		ret *= x
	}
	return ret
}
