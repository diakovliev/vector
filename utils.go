package vector

// NewStringsSet creates a new SetImpl[string, CompareFunc[string]] from the given string slice.
//
// data: A slice of strings to add to the set.
// returns: A pointer to the created SetImpl[string, CompareFunc[string]].
func NewStringsSet(data []string) (ret *SetImpl[string, CompareFunc[string]]) {
	ret = NewSet[string, CompareFunc[string]](CompareString[string])
	ret.Add(data...)
	return ret
}
