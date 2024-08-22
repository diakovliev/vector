package vector

import "strings"

// CompareFunc is a function that compares two values
type CompareFunc[T any] func(T, T) int

// Number is a types restriction interface for numbers
type Number interface {
	int | uint | int32 | uint32 | int64 | uint64 | float32 | float64
}

// CompareNumber compares two numbers of type N and returns -1 if the lhs is less than the rhs,
// 1 if the lhs is greater than the rhs, and 0 if they are equal.
//
// lhs: a number of type N
// rhs: a number of type N
// int: the result of the comparison (-1, 0, 1)
func CompareNumber[N Number](lhs N, rhs N) int {
	switch {
	case lhs < rhs:
		return -1
	case lhs > rhs:
		return 1
	default:
		return 0
	}
}

// String is a types restriction interface for strings
type String interface {
	string
}

// CompareString compares two strings lexicographically and returns an integer value.
//
// lhs: the first string to be compared.
// rhs: the second string to be compared.
//
// Returns an integer value: 0 if the strings are equal, -1 if lhs is less than rhs,
// and 1 if lhs is greater than rhs.
func CompareString[S String](lhs S, rhs S) int {
	return strings.Compare(string(lhs), string(rhs))
}
