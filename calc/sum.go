package calc

func FastSum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func SlowSum(s []int) int {
	if len(s) == 0 {
		return 0
	}

	return s[0] + SlowSum(s[1:])
}
