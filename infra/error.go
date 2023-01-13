package infra

import "github.com/pkg/errors"

var (
	UnknownError = errors.New("unknown error")
)

func Assert(b bool, elseError error) {
	if !b {
		panic(elseError)
	}
}

func Catch(e *error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			*e = err
		} else {
			*e = errors.Wrapf(UnknownError, "value :%+v", r)
		}
	}
}
