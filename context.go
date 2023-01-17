package dyad

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type cloneContext struct {
	options   cloneOptions
	writePath func(w io.Writer)
}

func (c cloneContext) WithPath(
	format string,
	args ...any,
) cloneContext {
	previous := c.writePath
	c.writePath = func(w io.Writer) {
		if previous != nil {
			previous(w)
		}

		for i, a := range args {
			if a, ok := a.(reflect.Type); ok {
				args[i] = renderTypeName(a)
			}
		}

		fmt.Fprintf(w, format, args...)
	}

	return c
}

func (c cloneContext) Error(
	format string,
	args ...any,
) error {
	message := &strings.Builder{}

	c.writePath(message)
	message.WriteString(": ")
	fmt.Fprintf(message, format, args...)

	return errors.New(message.String())
}

func renderTypeName(t reflect.Type) string {
	typeName := t.String()
	if typeName == "interface {}" {
		return "any"
	}

	return typeName
}
