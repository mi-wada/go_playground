package sum

func Sum(a []int) int {
	res := 0

	for _, v := range a {
		res += v
	}

	return res
}

func sumRec(a []int) int {
	if len(a) == 0 {
		return 0
	}

	return a[0] + sumRec(a[1:])
}

func SumParallel(a []int, parallelCount int) int {
	if parallelCount <= 1 {
		return Sum(a)
	}

	res := 0
	ch := make(chan int)
	partSize := len(a) / parallelCount

	for i := 0; i < parallelCount; i++ {
		start := i * partSize
		end := start + partSize
		if i == parallelCount-1 {
			end = len(a)
		}

		go func(start, end int) {
			ch <- Sum(a[start:end])
		}(start, end)
	}

	for i := 0; i < parallelCount; i++ {
		res += <-ch
	}

	return res
}

func SumHeuristic(a []int) int {
	threshold := 1_000_000

	if len(a) > threshold {
		return Sum(a)
	} else {
		return SumParallel(a, 2)
	}
}

func SumHeuristicInline(a []int) int {
	threshold := 1_000_000

	if len(a) > threshold {
		res := 0

		for _, v := range a {
			res += v
		}

		return res
	} else {
		res := 0
		ch := make(chan int)

		go func() {
			ch <- Sum(a[:len(a)/2])
		}()

		go func() {
			ch <- Sum(a[len(a)/2:])
		}()

		res = <-ch + <-ch

		return res
	}
}
