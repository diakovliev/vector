package vector

// LockerStub is a stub for Locker
type LockerStub struct {
}

// MakeLockerStub returns a LockerStub struct.
//
// This function takes no parameters.
// It returns a LockerStub struct.
func MakeLockerStub() LockerStub {
	return LockerStub{}
}

// NewLockerStub creates a new LockerStub instance and returns a pointer to it.
//
// No parameters.
// Returns a pointer to a LockerStub.
func NewLockerStub() *LockerStub {
	ret := MakeLockerStub()
	return &ret
}

// Lock locks the LockerStub.
//
// No parameters.
// No return values.
func (lw *LockerStub) Lock() {
}

// Unlock releases the lock held by the LockerStub instance.
//
// No parameters.
// No return values.
func (lw *LockerStub) Unlock() {
}
