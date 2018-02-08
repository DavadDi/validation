package validation

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// Error for Validater impl

	// ErrBadURLFormat for url
	ErrBadURLFormat = errors.New("url format is not valid")

	// ErrBadEmailFormat for email
	ErrBadEmailFormat = errors.New("email format is not valid")

	// ErrRequired for Required
	ErrRequired = errors.New("field can't be empty or zero")

	// Error for Validater

	// ErrValidater validater not exist
	ErrValidater = errors.New("validater should not be nil")

	// ErrValidaterNoFound can't find corresponding validater
	ErrValidaterNoFound = errors.New("validater not found")

	// ErrValidaterExists for registe validator repeatedly
	ErrValidaterExists = errors.New("validater exist")
)

// Error for Validator, including filedname, value, err msg.
type Error struct {
	FieldName string
	Value     interface{}
	Err       error
}

func (err *Error) String() string {
	return fmt.Sprintf("[%s] check failed [%s] [%#v]", err.FieldName, err.Err.Error(), err.Value)
}

// ErrUnsupportedType not support type
type ErrUnsupportedType struct {
	Type reflect.Type
}

func (err *ErrUnsupportedType) Error() string {
	return "validition unsupported type: " + err.Type.String()

}

// ErrOnlyStrcut only for validate struct
type ErrOnlyStrcut struct {
	Type reflect.Type
}

// ErrOnlyStrcut detail error msg
func (err *ErrOnlyStrcut) Error() string {
	return "validition only support struct, but got type: " + err.Type.String()

}

// NewErrWrongType new error for unmatched type
func NewErrWrongType(expect string, value interface{}) error {
	return &ErrWrongExpectType{
		ExpectType: expect,
		PassValue:  value,
	}
}

// ErrWrongExpectType  expect type don't match passed type
type ErrWrongExpectType struct {
	ExpectType string
	PassValue  interface{}
}

// ErrWrongExpectType detail error
func (err *ErrWrongExpectType) Error() string {
	return fmt.Sprintf("expect type %s, but got %T", err.ExpectType, err.PassValue)
}
