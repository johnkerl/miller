package types

// List is a thin wrapper around a slice which functions as a list/queue.  The Items attribute is
// exposed as public since I want iteration (at many, many callsites) to be easy.
type List[T any] struct {
	Items []T
}

func NewList[T any](capacity int) *List[T] {
	return &List[T]{
		make([]T, 0, capacity),
	}
}

// Front will panic if the list is empty
func (ell *List[T]) Front() T {
	return ell.Items[0]
}

func (ell *List[T]) Len() int {
	return len(ell.Items)
}

func (ell *List[T]) PushBack(e T) {
	ell.Items = append(ell.Items, e)
}

func (ell *List[T]) PushBackMultiple(mell []T) {
	ell.Items = append(ell.Items, mell...)
}

func (ell *List[T]) Clear() {
	ell.Items = ell.Items[0:0:cap(ell.Items)]
}
