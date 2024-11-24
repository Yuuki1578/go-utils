package vector

import (
	"errors"
	"slices"
)

// Compile-time constant errors
const (
	NIL_VALUE_ACCESS    string = "Error: attempting to access nil value"
	INDEX_OUT_OF_BOUNDS string = "Error: attempting to access invalid memory location"
)

// Basic vector, capable of appending, popping, removing, etc.
type Vector[T any] struct {
	__slice []T
	__len   uint64
	__cap   uint64
}

// Get the vector length
func (this *Vector[T]) Len() uint64 {
	if this != nil {
		return this.__len
	}

	return 0
}

// Get the vector capacity
func (this *Vector[T]) Cap() uint64 {
	if this != nil {
		return this.__cap
	}

	return 0
}

// Update the length and capacity of the vector.
func (this *Vector[T]) updateStatus() {
	if this.__slice == nil {
		return
	}

	this.__len = uint64(len(this.__slice))
	this.__cap = uint64(cap(this.__slice))
}

// Initialize a vector with default capacity and allocates the slices with the capacity provided first.
// In fact, this function is an abstraction over builtin function 'make'.
func WithCapacity[T any](capacity uint64) Vector[T] {
	return Vector[T]{
		__slice: make([]T, 0, capacity),
		__len:   0,
		__cap:   capacity,
	}
}

// Initialize a vector, it doesn't allocate the memory yet.
// The initial capacity is 0. If you want to specify the default
// capacity, use 'WithCapacity' instead.
func New[T any]() Vector[T] {
	return WithCapacity[T](0)
}

// Adding additional capacity to vector, if the pointer to instance (that called by the method)
// is somehow 'nil', it return error. If it's NOT nil, it return nil instead.
func (this *Vector[T]) AddCapacity(capacity uint64) error {
	var safeCopy []T

	if this == nil {
		return errors.New(NIL_VALUE_ACCESS)
	}

	if this.__slice == nil {
		*this = WithCapacity[T](capacity)
		return nil
	}

	var _ = copy(safeCopy, this.__slice)

	this.__slice = nil
	*this = WithCapacity[T](uint64(cap(safeCopy) + int(capacity)))
	this.__slice = append(this.__slice, safeCopy...)
	this.updateStatus()

	safeCopy = nil

	return nil
}

// Clear the vector, truncating it to initialization / zero value.
func (this *Vector[T]) Clear() {
	if this == nil {
		return
	}

	this.__slice = nil
	*this = New[T]()
	this.updateStatus()
}

// Deallocates the remaining capacity of vector.
func (this *Vector[T]) Strip() {
	if this == nil {
		return
	}

	this.__slice = slices.Clip(this.__slice)
	this.updateStatus()
}

// Reversing the vector
func (this *Vector[T]) Reverse() error {
	if this == nil || this.__slice == nil {
		return errors.New(NIL_VALUE_ACCESS)
	}

	slices.Reverse(this.__slice)

	return nil
}

// Appending element from the back, return error if the instance is nil.
// On success, return nil instead.
func (this *Vector[T]) Append(element T) error {
	if this == nil {
		return errors.New(NIL_VALUE_ACCESS)
	}

	if this.__slice == nil {
		*this = New[T]()
	}

	this.__slice = append(this.__slice, element)
	this.updateStatus()

	return nil
}

// Popping the value out of the vector, if the instance is nil or the index is greater than / equal to instance length,
// It will return the default value of 'T' and error.
func (this *Vector[T]) Pop(index uint64) (T, error) {
	var (
		defaultValue T
		leftSide     []T = nil
		rightSide    []T = nil
		safeCopy     []T = nil
		_            int
	)

	if this == nil || this.__slice == nil {
		return defaultValue, errors.New(NIL_VALUE_ACCESS)
	}

	if this.Len() >= index {
		return defaultValue, errors.New(INDEX_OUT_OF_BOUNDS)
	}

	defaultValue = this.__slice[index]
	leftSide = this.__slice[:index]
	rightSide = this.__slice[index+1:]

	safeCopy = make([]T, 0, this.Cap()-1)
	safeCopy = append(safeCopy, leftSide...)
	safeCopy = append(safeCopy, rightSide...)

	this.Clear()

	_ = copy(this.__slice, safeCopy)
	this.updateStatus()

	safeCopy = nil

	return defaultValue, nil
}
