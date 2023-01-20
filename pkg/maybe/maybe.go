package maybe

import "errors"

type BindFunc[A, B any] func(A) (B, error)

type Maybe[A any] interface {
	Error() error
	ValueE() (A, error)
	Value() A
}

func Bind[A, B any](a Maybe[A], bindFunc BindFunc[A, B]) Maybe[B] {
	if a.Error() != nil {
		return NewFailure[B](a.Error())
	}

	return Result(bindFunc(a.Value()))
}

type Success[A any] struct {
	value A
}

type Failure[A any] struct {
	err error
}

func (success *Success[A]) Error() error {
	return nil
}

func (success *Success[A]) ValueE() (A, error) {
	return success.value, nil
}

func (success *Success[A]) Value() A {
	return success.value
}

var ErrorNoValue = errors.New("cannot get value from failure Maybe")

func (failure *Failure[A]) Error() error {
	return failure.err
}

func (failure *Failure[A]) ValueE() (A, error) {
	return *new(A), ErrorNoValue
}

func (failure *Failure[A]) Value() A {
	return *new(A)
}

func NewSuccess[A any](value A) *Success[A] {
	return &Success[A]{
		value: value,
	}
}

func NewFailure[A any](err error) *Failure[A] {
	return &Failure[A]{
		err: err,
	}
}

func Result[A any](val A, err error) Maybe[A] {
	if err != nil {
		return NewFailure[A](err)
	}

	return NewSuccess(val)
}
