package errors

import "fmt"

type MultiError struct {
	Errors []error
}

func (me MultiError) Error() string {
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors, %d others and: %s", len(me.Errors), me.Errors[0].Error())
}

