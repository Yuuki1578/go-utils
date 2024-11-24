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
	Slice []T
	Len   uint64
	Cap   uint64
}

// Update the length and capacity of the vector.
func (this *Vector[T]) updateStatus() {
	if this == nil || this.Slice == nil {
		return
	}

	this.Len = uint64(len(this.Slice))
	this.Cap = uint64(cap(this.Slice))
}

// Initialize a vector with default capacity and allocates the slices with the capacity provided first.
// In fact, this function is an abstraction over builtin function 'make'.
func WithCapacity[T any](capacity uint64) Vector[T] {
	return Vector[T]{
		Slice: make([]T, 0, capacity),
		Len:   0,
		Cap:   capacity,
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

	if this.Slice == nil {
		*this = WithCapacity[T](capacity)
		return nil
	}

	var _ = copy(safeCopy, this.Slice)

	this.Slice = nil
	*this = WithCapacity[T](uint64(cap(safeCopy) + int(capacity)))
	this.Slice = append(this.Slice, safeCopy...)
	this.updateStatus()

	safeCopy = nil

	return nil
}

// Clear the vector, truncating it to initialization / zero value.
func (this *Vector[T]) Clear() {
	if this == nil {
		return
	}

	this.Slice = nil
	*this = New[T]()
	this.updateStatus()
}

// Deallocates the remaining capacity of vector.
func (this *Vector[T]) Strip() {
	if this == nil {
		return
	}

	this.Slice = slices.Clip(this.Slice)
	this.updateStatus()
}

// Reversing the vector
func (this *Vector[T]) Reverse() error {
	if this == nil || this.Slice == nil {
		return errors.New(NIL_VALUE_ACCESS)
	}

	slices.Reverse(this.Slice)

	return nil
}

// Appending element from the back, return error if the instance is nil.
// On success, return nil instead.
func (this *Vector[T]) Append(element T) error {
	if this == nil {
		return errors.New(NIL_VALUE_ACCESS)
	}

	if this.Slice == nil {
		*this = New[T]()
	}

	this.Slice = append(this.Slice, element)
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

	if this == nil || this.Slice == nil {
		return defaultValue, errors.New(NIL_VALUE_ACCESS)
	}

	if this.Len >= index {
		return defaultValue, errors.New(INDEX_OUT_OF_BOUNDS)
	}

	defaultValue = this.Slice[index]
	leftSide = this.Slice[:index]
	rightSide = this.Slice[index+1:]

	safeCopy = make([]T, 0, this.Cap-1)
	safeCopy = append(safeCopy, leftSide...)
	safeCopy = append(safeCopy, rightSide...)

	this.Clear()

	_ = copy(this.Slice, safeCopy)
	this.updateStatus()

	safeCopy = nil

	return defaultValue, nil
}
