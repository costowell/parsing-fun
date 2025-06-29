package common

import "fmt"

// OrderedSet is effectively an array with guaranteed unique elements
type OrderedSet[T comparable] struct {
	indexMap map[T]int
	Data     []T
}

// String returns the string representation of the underlying data
func (o *OrderedSet[T]) String() string {
	return fmt.Sprintf("%+v", o.Data)
}

// NewOrderedSet creates a new OrderedSet
func NewOrderedSet[T comparable]() OrderedSet[T] {
	return OrderedSet[T]{
		indexMap: make(map[T]int),
		Data:     make([]T, 0),
	}
}

// Insert adds an element to the end of the set, returns false if the element already exists
func (o *OrderedSet[T]) Insert(elm T) bool {
	if o.Contains(elm) {
		return false
	}
	o.indexMap[elm] = len(o.Data)
	o.Data = append(o.Data, elm)
	return true
}

// Contains returns whether an element already exists in the set
func (o *OrderedSet[T]) Contains(elm T) bool {
	_, exists := o.indexMap[elm]
	return exists
}

// Remove deletes an element from the set, returns false if it doesn't exist
func (o *OrderedSet[T]) Remove(elm T) bool {
	index, exists := o.indexMap[elm]
	if exists {
		o.Data = append(o.Data[:index], o.Data[index:]...)
		delete(o.indexMap, elm)
	}
	return exists
}
