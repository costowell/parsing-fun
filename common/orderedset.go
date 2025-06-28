package common

import "fmt"

type OrderedSet[T comparable] struct {
	existMap map[T]bool
	Data     []T
}

func (o *OrderedSet[T]) String() string {
	return fmt.Sprintf("%+v", o.Data)
}

func NewOrderedSet[T comparable]() OrderedSet[T] {
	return OrderedSet[T]{
		existMap: make(map[T]bool),
		Data:     make([]T, 0),
	}
}

func (o *OrderedSet[T]) Insert(elm T) bool {
	if o.Exists(elm) {
		return false
	}
	o.existMap[elm] = true
	o.Data = append(o.Data, elm)
	return true
}

func (o *OrderedSet[T]) Exists(elm T) bool {
	_, exists := o.existMap[elm]
	return exists
}
