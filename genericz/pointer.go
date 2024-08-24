package genericz

// Pointer returns a pointer to the value passed.
//
// Pointer and Ptr are the same.
func Pointer[T any](v T) *T { return &v }

// Ptr returns a pointer to the value passed.
//
// Pointer and Ptr are the same.
func Ptr[T any](v T) *T { return &v }
