package pago

func min[T int | uint](a, b T) T {
	if a < b {
		return a
	}
	return b
}
