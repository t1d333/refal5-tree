package runtime

func fibonacci() func(n int) int {
	memo := map[int]int{}

	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if value, ok := memo[n]; ok {
			return value
		}
		memo[n] = fib(n-1) + fib(n-2)
		return memo[n]
	}

	return fib
}
