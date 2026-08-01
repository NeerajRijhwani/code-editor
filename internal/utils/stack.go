package utils

import (
	"errors"
)

type Stack[T any] []T

func (s *Stack[T]) Push(val T) {
	*s = append(*s, val)
}

func Init_Stack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Pop() (val T, err error) {
	if s.isEmpty() {
		var zero T
		return zero, errors.New("Stack is Empty")
	}
	top := len(*s) - 1

	val = (*s)[top]

	var zero T
	(*s)[top] = zero
	*s = (*s)[:top]

	return val, nil

}

func (s *Stack[T]) isEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Clear() {
	var zero T
	for i := range *s {
		(*s)[i] = zero
	}
	*s = (*s)[:0]
}
